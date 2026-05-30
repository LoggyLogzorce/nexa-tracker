package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/models"
	"nexa-task-tracker/internal/pkg/hash"
	jwtpkg "nexa-task-tracker/internal/pkg/jwt"
	"time"
)

type GoogleOAuthService interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string, userAgent, ip *string) (string, error)
}

type googleOAuthService struct {
	config        *oauth2.Config
	authRepo      auth.Repository
	userRepo      user.Repository
	jwtSecret     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewGoogleOAuthService(config *oauth2.Config, authRepo auth.Repository, userRepo user.Repository,
	secret string, accessExp, refreshExp time.Duration) GoogleOAuthService {
	return &googleOAuthService{
		config:        config,
		authRepo:      authRepo,
		userRepo:      userRepo,
		jwtSecret:     secret,
		AccessExpiry:  accessExp,
		RefreshExpiry: refreshExp,
	}
}

type GoogleUser struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

func (s *googleOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *googleOAuthService) Exchange(ctx context.Context, code string, userAgent, ip *string) (string, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, err := s.config.Exchange(ctxT, code)
	if err != nil {
		return "", err
	}

	client := s.config.Client(ctxT, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var googleUser GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return "", err
	}

	// ищем провайдера
	provider, err := s.authRepo.GetByProviderAndID(ctxT, GoogleProvider, googleUser.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	var u *models.User

	if provider.ID != 0 {
		// провайдер есть — просто берём юзера
		u, err = s.userRepo.GetByID(ctxT, provider.UserID)
		if err != nil {
			return "", err
		}
	} else {
		// ищем юзера по email
		u, err = s.userRepo.GetByEmail(ctxT, googleUser.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || u == nil {
			// создаём нового юзера
			u = &models.User{
				ID:           uuid.New(),
				Email:        googleUser.Email,
				Name:         googleUser.Name,
				AvatarUrl:    googleUser.Picture,
				PasswordHash: "",
				Role:         "user",
			}
			if err := s.userRepo.Create(ctxT, u); err != nil {
				return "", err
			}
		}

		// создаём провайдера
		up := &models.UserProvider{
			UserID:         u.ID,
			Provider:       GoogleProvider,
			ProviderUserID: googleUser.ID,
		}
		if err := s.authRepo.CreateUserProvider(ctxT, up); err != nil {
			return "", err
		}
	}

	return s.issueTokens(ctxT, u, userAgent, ip)
}

func (s *googleOAuthService) issueTokens(ctx context.Context, u *models.User, userAgent, ip *string) (string, error) {
	refreshToken, err := jwtpkg.GenerateRefreshToken(u.ID, s.jwtSecret, s.RefreshExpiry)
	if err != nil {
		return "", err
	}

	tokenHash := hash.TokenHash(refreshToken)

	rt := &models.RefreshToken{
		UserID:    u.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(s.RefreshExpiry),
		UserAgent: userAgent,
		IPAddress: ip,
	}
	if err := s.authRepo.CreateRefreshToken(ctx, rt); err != nil {
		return "", err
	}

	return refreshToken, nil
}

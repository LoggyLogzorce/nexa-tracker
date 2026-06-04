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
	"nexa-task-tracker/pkg/hash"
	jwtpkg "nexa-task-tracker/pkg/jwt"
	"time"
)

type YandexOAuthService interface {
	GetAuthURL(state string) string
	Exchange(ctx context.Context, code string, userAgent, ip *string) (string, error)
}

type yandexOAuthService struct {
	config        *oauth2.Config
	authRepo      auth.Repository
	userRepo      user.Repository
	jwtSecret     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewYandexOAuthService(config *oauth2.Config, authRepo auth.Repository, userRepo user.Repository,
	secret string, accessExp, refreshExp time.Duration) YandexOAuthService {
	return &yandexOAuthService{
		config:        config,
		authRepo:      authRepo,
		userRepo:      userRepo,
		jwtSecret:     secret,
		AccessExpiry:  accessExp,
		RefreshExpiry: refreshExp,
	}
}

type YandexUser struct {
	ID              string `json:"id"`
	Login           string `json:"login"`
	RealName        string `json:"real_name"`
	DefaultEmail    string `json:"default_email"`
	IsAvatarEmpty   bool   `json:"is_avatar_empty"`
	DefaultAvatarID string `json:"default_avatar_id"`
}

func (u *YandexUser) AvatarURL() string {
	if u.IsAvatarEmpty || u.DefaultAvatarID == "" {
		return ""
	}
	return "https://avatars.yandex.net/get-yapic/" + u.DefaultAvatarID + "/islands-200"
}

func (s *yandexOAuthService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *yandexOAuthService) Exchange(ctx context.Context, code string, userAgent, ip *string) (string, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	token, err := s.config.Exchange(ctxT, code)
	if err != nil {
		return "", err
	}

	client := s.config.Client(ctxT, token)
	resp, err := client.Get("https://login.yandex.ru/info?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var yandexUser YandexUser
	if err := json.NewDecoder(resp.Body).Decode(&yandexUser); err != nil {
		return "", err
	}

	provider, err := s.authRepo.GetByProviderAndID(ctxT, YandexProvider, yandexUser.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	var u *models.User

	if provider.ID != 0 {
		u, err = s.userRepo.GetByID(ctxT, provider.UserID)
		if err != nil {
			return "", err
		}
	} else {
		u, err = s.userRepo.GetByEmail(ctxT, yandexUser.DefaultEmail)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || u == nil {
			u = &models.User{
				ID:           uuid.New(),
				Email:        yandexUser.DefaultEmail,
				Name:         yandexUser.RealName,
				AvatarUrl:    yandexUser.AvatarURL(),
				PasswordHash: "",
				Role:         "user",
			}
			if err := s.userRepo.Create(ctxT, u); err != nil {
				return "", err
			}
		}

		up := &models.UserProvider{
			UserID:         u.ID,
			Provider:       YandexProvider,
			ProviderUserID: yandexUser.ID,
		}
		if err := s.authRepo.CreateUserProvider(ctxT, up); err != nil {
			return "", err
		}
	}

	return s.issueTokens(ctxT, u, userAgent, ip)
}

func (s *yandexOAuthService) issueTokens(ctx context.Context, u *models.User, userAgent, ip *string) (string, error) {
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

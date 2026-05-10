package auth

import (
	"context"
	"errors"
	"time"

	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/pkg/hash"
	jwtpkg "nexa-task-tracker/internal/pkg/jwt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	Register(ctx context.Context, email, password string, name *string) error
	Login(ctx context.Context, email, password, userAgent, ipAddress string) (accessToken, refreshToken string, err error)
	RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error)
	Logout(ctx context.Context, refreshToken string) error
	Setup2FA(ctx context.Context, userID string) (secret, qrCode string, err error)
	Verify2FA(ctx context.Context, userID, code string) error
	Enable2FA(ctx context.Context, userID, code string) error
	Disable2FA(ctx context.Context, userID, code string) error
}

type service struct {
	repo          Repository
	userRepo      user.Repository
	jwtSecret     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

func NewService(repo Repository, userRepo user.Repository, jwtSecret string, accessExpiry, refreshExpiry time.Duration) Service {
	return &service{
		repo:          repo,
		userRepo:      userRepo,
		jwtSecret:     jwtSecret,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}
}

func (s *service) Register(ctx context.Context, email, password string, name *string) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if email already exists
	_, err := s.userRepo.GetByEmail(ctxT, email)
	if err == nil {
		// User found, email already exists
		return ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// Database error
		return err
	}

	// Hash password
	hashedPassword, err := hash.Generate(password)
	if err != nil {
		return err
	}

	// Create user
	userName := ""
	if name != nil {
		userName = *name
	}

	newUser := &user.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         userName,
		Role:         "user",
	}

	return s.userRepo.Create(ctxT, newUser)
}

func (s *service) Login(ctx context.Context, email, password, userAgent, ipAddress string) (accessToken, refreshToken string, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Find user by email
	foundUser, err := s.userRepo.GetByEmail(ctxT, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	// Verify password
	if err := hash.Compare(foundUser.PasswordHash, password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	// Generate access token (15 minutes default)
	accessToken, err = jwtpkg.GenerateAccessToken(
		foundUser.ID,
		foundUser.Email,
		foundUser.Role,
		s.jwtSecret,
		s.AccessExpiry,
	)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (7 days default)
	refreshToken, err = jwtpkg.GenerateRefreshToken(
		foundUser.ID,
		s.jwtSecret,
		s.RefreshExpiry,
	)
	if err != nil {
		return "", "", err
	}

	// Hash refresh token for storage
	hashedRefreshToken := hash.TokenHash(refreshToken)

	// Save refresh token to database
	refreshTokenRecord := &RefreshToken{
		UserID:    foundUser.ID,
		TokenHash: hashedRefreshToken,
		ExpiresAt: time.Now().Add(s.RefreshExpiry),
		UserAgent: &userAgent,
		IPAddress: &ipAddress,
	}

	if err := s.repo.CreateRefreshToken(ctxT, refreshTokenRecord); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (accessToken, newRefreshToken string, err error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Валидировать refresh token JWT
	claims, err := jwtpkg.Validate(refreshToken, s.jwtSecret)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// 2. Хешировать refresh token для поиска в БД
	tokenHash := hash.TokenHash(refreshToken)

	// 3. Найти refresh token в БД
	tokenRecord, err := s.repo.GetRefreshToken(ctxT, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", ErrInvalidToken
		}
		return "", "", err
	}

	// 4. Проверить revocation (reuse attack detection)
	if tokenRecord.RevokedAt != nil {
		// Токен уже использован! Подозрение на атаку
		// Отозвать все токены пользователя
		s.repo.RevokeAllUserTokens(ctxT, tokenRecord.UserID)
		return "", "", ErrInvalidToken
	}

	// 5. Проверить expiration
	if tokenRecord.ExpiresAt.Before(time.Now()) {
		// Токен истёк, удалить из БД
		s.repo.DeleteExpiredTokens(ctxT)
		return "", "", ErrInvalidToken
	}

	// 6. Отозвать текущий refresh token (rotation)
	if err := s.repo.RevokeRefreshToken(ctxT, tokenHash); err != nil {
		return "", "", err
	}

	// 7. Найти пользователя в БД
	foundUser, err := s.userRepo.GetByID(ctxT, claims.UserID)
	if err != nil {
		return "", "", ErrInvalidToken
	}

	// 8. Сгенерировать новый access token
	newAccessToken, err := jwtpkg.GenerateAccessToken(
		foundUser.ID,
		foundUser.Email,
		foundUser.Role,
		s.jwtSecret,
		s.AccessExpiry,
	)
	if err != nil {
		return "", "", err
	}

	// 9. Сгенерировать новый refresh token
	newRefreshToken, err = jwtpkg.GenerateRefreshToken(
		foundUser.ID,
		s.jwtSecret,
		s.RefreshExpiry, // 7 days
	)
	if err != nil {
		return "", "", err
	}

	// 10. Хешировать новый refresh token
	newTokenHash := hash.TokenHash(newRefreshToken)

	// 11. Сохранить новый refresh token в БД
	newTokenRecord := &RefreshToken{
		UserID:    foundUser.ID,
		TokenHash: newTokenHash,
		ExpiresAt: time.Now().Add(168 * time.Hour),
		UserAgent: tokenRecord.UserAgent, // сохранить из старой записи
		IPAddress: tokenRecord.IPAddress,
	}

	if err := s.repo.CreateRefreshToken(ctxT, newTokenRecord); err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *service) Logout(ctx context.Context, refreshToken string) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Хешировать refresh token для поиска в БД
	tokenHash := hash.TokenHash(refreshToken)

	// Отозвать refresh token в БДctxT,
	return s.repo.RevokeRefreshToken(ctxT, tokenHash)
}

func (s *service) Setup2FA(ctx context.Context, userID string) (secret, qrCode string, err error) {
	// TODO: Implement TOTP setup
	return "", "", nil
}

func (s *service) Verify2FA(ctx context.Context, userID, code string) error {
	// TODO: Implement TOTP verification
	return nil
}

func (s *service) Enable2FA(ctx context.Context, userID, code string) error {
	// TODO: Implement
	return nil
}

func (s *service) Disable2FA(ctx context.Context, userID, code string) error {
	// TODO: Implement
	return nil
}

package auth

import (
	"context"
	"nexa-task-tracker/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error
	GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string) error
	RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error
	DeleteExpiredTokens(ctx context.Context) error
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.RefreshToken, error)
	DeleteSessionByID(ctx context.Context, userID uuid.UUID, sessionID uint) error

	GetByProviderAndID(ctx context.Context, provider string, id string) (*models.UserProvider, error)
	CreateUserProvider(ctx context.Context, up *models.UserProvider) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *repository) GetRefreshToken(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := r.db.WithContext(ctx).Where("token_hash = ?", tokenHash).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", now).Error
}

func (r *repository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func (r *repository) DeleteExpiredTokens(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

func (r *repository) GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.RefreshToken, error) {
	var tokens []models.RefreshToken
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND revoked_at IS NULL AND expires_at > ?", userID, time.Now()).
		Order("created_at DESC").
		Find(&tokens).Error
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *repository) DeleteSessionByID(ctx context.Context, userID uuid.UUID, sessionID uint) error {
	return r.db.WithContext(ctx).Model(&models.RefreshToken{}).
		Where("id = ? AND user_id = ? AND revoked_at IS NULL", sessionID, userID).
		Delete(&models.RefreshToken{}).Error
}

func (r *repository) GetByProviderAndID(ctx context.Context, provider string, id string) (*models.UserProvider, error) {
	var up models.UserProvider
	err := r.db.WithContext(ctx).Where(&models.UserProvider{Provider: provider, ProviderUserID: id}).First(&up).Error
	return &up, err
}

func (r *repository) CreateUserProvider(ctx context.Context, up *models.UserProvider) error {
	return r.db.WithContext(ctx).Create(up).Error
}

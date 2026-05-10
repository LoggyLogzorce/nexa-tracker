package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/pkg/events"
	"nexa-task-tracker/internal/pkg/hash"
)

type Service interface {
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID, password string) error
	EmailExists(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error)
}

type service struct {
	repo     Repository
	eventBus *events.EventBus
}

func NewService(repo Repository, eventBus *events.EventBus) Service {
	return &service{
		repo:     repo,
		eventBus: eventBus,
	}
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.GetByID(ctxT, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (*User, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(ctx context.Context, user *User) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Update(ctxT, user)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, password string) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить пользователя
	user, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Проверить пароль
	if err := hash.Compare(user.PasswordHash, password); err != nil {
		return ErrInvalidPassword
	}

	// 3. Проверить что пользователь не владелец проектов
	hasProjects, err := s.repo.UserOwnsProjects(ctxT, id)
	if err != nil {
		return fmt.Errorf("failed to check projects: %w", err)
	}
	if hasProjects {
		return ErrUserOwnsProjects
	}

	// 4. Анонимизировать данные перед удалением
	user.Email = fmt.Sprintf("deleted_%s@deleted.local", id.String())
	user.Name = "Deleted User"
	user.PasswordHash = ""
	user.Secret2FA = nil

	// 5. Сохранить анонимизированные данные
	if err := s.repo.Update(ctxT, user); err != nil {
		return fmt.Errorf("failed to anonymize user: %w", err)
	}

	//// 6. Soft delete пользователя
	//if err := s.repo.Delete(id); err != nil {
	//	return fmt.Errorf("failed to delete user: %w", err)
	//}

	// 7. Опубликовать событие UserDeleted
	event := UserDeletedEvent{
		UserID: id,
		Email:  user.Email,
		Name:   user.Name,
	}
	s.eventBus.Publish(event.ToEvent())

	return nil
}

func (s *service) EmailExists(ctx context.Context, email string, excludeUserID uuid.UUID) (bool, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Привести email к lowercase для case-insensitive сравнения
	email = strings.ToLower(email)

	user, err := s.repo.GetByEmail(ctxT, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // Email свободен
		}
		return false, err // Ошибка БД
	}

	// Email занят другим пользователем?
	return user.ID != excludeUserID, nil
}

package user

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	SearchByEmail(ctx context.Context, query string) ([]UserResponse, error)
	Update(ctx context.Context, user *User) (*UserResponse, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, filename string, file io.Reader, uploadPath string) (*UserResponse, error)
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

func (s *service) SearchByEmail(ctx context.Context, query string) ([]UserResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	users, err := s.repo.SearchByEmail(ctxT, query)
	if err != nil {
		return nil, err
	}

	responses := make([]UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, *users[i].ToResponse())
	}
	return responses, nil
}

func (s *service) Update(ctx context.Context, user *User) (*UserResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	userOld, err := s.repo.GetByID(ctxT, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	userOld.Name = user.Name

	if err := s.repo.Update(ctxT, userOld); err != nil {
		return nil, err
	}
	return userOld.ToResponse(), nil
}

func (s *service) UploadAvatar(ctx context.Context, userID uuid.UUID, filename string, file io.Reader, uploadPath string) (*UserResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	user, err := s.repo.GetByID(ctxT, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	fileUID := uuid.New().String()
	ext := filepath.Ext(filename)
	storedName := fileUID + ext

	avatarDir := filepath.Join(uploadPath, "avatars")
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		log.Printf("Failed to create avatar directory: %v", err)
		return nil, fmt.Errorf("failed to create avatar directory: %w", err)
	}

	filePath := filepath.Join(avatarDir, storedName)
	dst, err := os.Create(filePath)
	if err != nil {
		log.Printf("Failed to create avatar file: %v", err)
		return nil, fmt.Errorf("failed to create avatar file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		os.Remove(filePath)
		log.Printf("Failed to write avatar file: %v", err)
		return nil, fmt.Errorf("failed to write avatar file: %w", err)
	}

	// Remove old avatar if exists
	if user.AvatarUrl != "" {
		if err := os.Remove(user.AvatarUrl); err != nil && !os.IsNotExist(err) {
			log.Printf("Failed to remove old avatar: %v", err)
		}
	}

	user.AvatarUrl = "/uploads/avatars/" + storedName
	if err := s.repo.Update(ctxT, user); err != nil {
		os.Remove(filePath)
		return nil, err
	}

	return user.ToResponse(), nil
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

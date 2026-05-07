package participant

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Service interface {
	AddParticipant(ctx context.Context, participant *ProjectParticipant) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]ProjectParticipant, error)
	GetByUserID(ctx context.Context, userID string) ([]ProjectParticipant, error)
	UpdateRole(ctx context.Context, participant *ProjectParticipant) error
	RemoveParticipant(ctx context.Context, participant *ProjectParticipant) error
	CheckAccess(ctx context.Context, projectID uuid.UUID, userID string, requiredRole string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddParticipant(ctx context.Context, participant *ProjectParticipant) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := s.repo.GetByProjectAndUser(ctxT, participant.ProjectID, participant.UserID)
	if err == nil {
		return ErrParticipantIDExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if err := s.repo.Create(ctxT, participant); err != nil {
		return err
	}

	return nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]ProjectParticipant, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.GetByProjectID(ctxT, projectID)
}

func (s *service) GetByUserID(ctx context.Context, userID string) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) UpdateRole(ctx context.Context, participant *ProjectParticipant) error {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.repo.GetByProjectAndUser(ctxT, participant.ProjectID, participant.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrParticipantNotFound
		}
		return err
	}

	return s.repo.Update(ctxT, participant)
}

func (s *service) RemoveParticipant(ctx context.Context, participant *ProjectParticipant) error {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.repo.GetByProjectAndUser(ctxT, participant.ProjectID, participant.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrParticipantNotFound
		}
		return err
	}

	return s.repo.Delete(ctxT, participant)
}

func (s *service) CheckAccess(ctx context.Context, projectID uuid.UUID, userID string, requiredRole string) (bool, error) {
	// TODO: Implement role hierarchy check
	return false, nil
}

// TODO удаление участников при удалении проекта

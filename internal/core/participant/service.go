package participant

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/user"
	"time"
)

type Service interface {
	AddParticipant(ctx context.Context, participant *ProjectParticipant) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]ProjectParticipantsResponse, error)
	GetByUserID(ctx context.Context, userID string) ([]ProjectParticipant, error)
	UpdateRole(ctx context.Context, participant *ProjectParticipant) error
	RemoveParticipant(ctx context.Context, participant *ProjectParticipant) error
	CheckAccess(ctx context.Context, projectID uuid.UUID, userID string, requiredRole string) (bool, error)
}

type service struct {
	repo     Repository
	userRepo user.Repository
}

func NewService(repo Repository, userRepo user.Repository) Service {
	return &service{
		repo:     repo,
		userRepo: userRepo,
	}
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

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]ProjectParticipantsResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	participants, err := s.repo.GetByProjectID(ctxT, projectID)
	if err != nil {
		return nil, err
	}

	participantsResponse := make([]ProjectParticipantsResponse, len(participants))
	userIds := make([]uuid.UUID, len(participants))
	for e, v := range participants {
		participantsResponse[e] = ProjectParticipantsResponse{
			ProjectID: v.ProjectID,
			Role:      v.Role,
			User: struct {
				UserID uuid.UUID `json:"user_id"`
				Name   string    `json:"name"`
				Email  string    `json:"email"`
			}{UserID: v.UserID, Name: "", Email: ""},
		}
		userIds[e] = v.UserID
	}

	users, err := s.userRepo.GetListByIDs(ctxT, userIds)
	if err != nil {
		return nil, err
	}

	for i := range participantsResponse {
		for _, u := range users {
			if participantsResponse[i].User.UserID == u.ID {
				participantsResponse[i].User.Name = u.Name
				participantsResponse[i].User.Email = u.Email
			}
		}
	}

	return participantsResponse, nil
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

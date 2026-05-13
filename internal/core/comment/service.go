package comment

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/user"
	"time"
)

type Service interface {
	Create(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, id uint) (*Comment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]CommentResponse, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id, taskID uint, userID uuid.UUID) error
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

func (s *service) Create(ctx context.Context, comment *Comment) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment.CreatedAt = time.Now()

	return s.repo.Create(ctxT, comment)
}

func (s *service) GetByID(ctx context.Context, id uint) (*Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByTaskID(ctx context.Context, taskID uint) ([]CommentResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comments, err := s.repo.GetByTaskID(ctxT, taskID)
	if err != nil {
		return nil, err
	}

	if len(comments) == 0 {
		return []CommentResponse{}, nil
	}

	userIDsMap := make(map[uuid.UUID]struct{})
	for _, c := range comments {
		userIDsMap[c.UserID] = struct{}{}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]user.User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	response := make([]CommentResponse, len(comments))
	for i, comment := range comments {
		response[i] = CommentResponse{
			ID:        comment.ID,
			CreatedAt: comment.CreatedAt,
			User:      CommentUserResponse{},
			TaskID:    comment.TaskID,
			Content:   comment.Content,
		}
		if u, ok := usersMap[comment.UserID]; ok {
			response[i].User = CommentUserResponse{
				ID:    u.ID,
				Name:  u.Name,
				Email: u.Email,
			}
		}
	}

	return response, nil
}

func (s *service) Update(ctx context.Context, comment *Comment) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	commentOld, err := s.repo.GetByID(ctxT, comment.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrCommentNotFound
	}

	if commentOld.TaskID != comment.TaskID {
		return ErrCommentNotFound
	}
	if commentOld.UserID != comment.UserID {
		return ErrNotCommentOwner
	}

	comment.ID = commentOld.ID
	comment.CreatedAt = commentOld.CreatedAt

	return s.repo.Update(ctxT, comment)
}

func (s *service) Delete(ctx context.Context, id, taskID uint, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment, err := s.repo.GetByID(ctxT, id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if comment.TaskID != taskID {
		return ErrCommentNotFound
	}

	if comment.UserID != userID {
		return ErrNotCommentOwner
	}

	return s.repo.Delete(ctxT, id)
}

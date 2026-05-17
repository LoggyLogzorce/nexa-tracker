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
	Create(ctx context.Context, comment *Comment) (*CommentResponse, error)
	GetByID(ctx context.Context, id uint) (*Comment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]CommentResponse, error)
	Update(ctx context.Context, comment *Comment) (*CommentResponse, error)
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

func (s *service) Create(ctx context.Context, comment *Comment) (*CommentResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	comment.CreatedAt = time.Now()

	if err := s.repo.Create(ctxT, comment); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctxT, comment.UserID)
	if err != nil {
		return nil, err
	}

	commentRes := &CommentResponse{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		User: CommentUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		},
		TaskID:  comment.TaskID,
		Content: comment.Content,
	}

	return commentRes, err
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
				ID:        u.ID,
				Name:      u.Name,
				Email:     u.Email,
				AvatarUrl: u.AvatarUrl,
			}
		}
	}

	return response, nil
}

func (s *service) Update(ctx context.Context, comment *Comment) (*CommentResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	commentOld, err := s.repo.GetByID(ctxT, comment.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrCommentNotFound
	}

	if commentOld.TaskID != comment.TaskID {
		return nil, ErrCommentNotFound
	}
	if commentOld.UserID != comment.UserID {
		return nil, ErrNotCommentOwner
	}

	comment.ID = commentOld.ID
	comment.CreatedAt = commentOld.CreatedAt

	if err := s.repo.Update(ctxT, comment); err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctxT, comment.UserID)
	if err != nil {
		return nil, err
	}

	commentRes := &CommentResponse{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		User: CommentUserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			AvatarUrl: user.AvatarUrl,
		},
		TaskID:  comment.TaskID,
		Content: comment.Content,
	}

	return commentRes, err
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

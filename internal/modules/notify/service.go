package notify

import "nexa-task-tracker/internal/pkg/events"

type Service interface {
	Create(notification *Notification) error
	GetByUserID(userID string, limit int) ([]Notification, error)
	GetUnreadByUserID(userID string) ([]Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(userID string) error
	NotifyTaskAssigned(taskID uint, assigneeID, assignerID string) error
	NotifyTaskUpdated(taskID uint, userID string) error
	NotifyCommentAdded(taskID uint, commentID uint, userID string) error
	HandleUserDeleted(event events.Event) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(notification *Notification) error {
	// TODO: Implement
	return nil
}

func (s *service) GetByUserID(userID string, limit int) ([]Notification, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetUnreadByUserID(userID string) ([]Notification, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) MarkAsRead(id uint) error {
	// TODO: Implement
	return nil
}

func (s *service) MarkAllAsRead(userID string) error {
	// TODO: Implement
	return nil
}

func (s *service) NotifyTaskAssigned(taskID uint, assigneeID, assignerID string) error {
	// TODO: Create notification for task assignment
	return nil
}

func (s *service) NotifyTaskUpdated(taskID uint, userID string) error {
	// TODO: Create notification for task update
	return nil
}

func (s *service) NotifyCommentAdded(taskID uint, commentID uint, userID string) error {
	// TODO: Create notification for new comment
	return nil
}

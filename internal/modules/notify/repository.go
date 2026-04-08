package notify

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(notification *Notification) error
	GetByID(id uint) (*Notification, error)
	GetByUserID(userID string, limit int) ([]Notification, error)
	GetUnreadByUserID(userID string) ([]Notification, error)
	MarkAsRead(id uint) error
	MarkAllAsRead(userID string) error
	DeleteByUserID(userID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(notification *Notification) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByID(id uint) (*Notification, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByUserID(userID string, limit int) ([]Notification, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetUnreadByUserID(userID string) ([]Notification, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) MarkAsRead(id uint) error {
	// TODO: Implement
	return nil
}

func (r *repository) MarkAllAsRead(userID string) error {
	// TODO: Implement
	return nil
}

func (r *repository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&Notification{}).Error
}

package comment

import (
	"context"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	GetByID(ctx context.Context, id uint) (*Comment, error)
	GetByTaskID(ctx context.Context, taskID uint) ([]Comment, error)
	Update(ctx context.Context, comment *Comment) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, comment *Comment) error {
	return r.db.WithContext(ctx).Create(comment).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Comment, error) {
	var comment Comment
	err := r.db.WithContext(ctx).First(&comment, id).Error
	return &comment, err
}

func (r *repository) GetByTaskID(ctx context.Context, taskID uint) ([]Comment, error) {
	var comments []Comment
	err := r.db.WithContext(ctx).Find(&comments, "task_id = ?", taskID).Error
	return comments, err
}

func (r *repository) Update(ctx context.Context, comment *Comment) error {
	return r.db.WithContext(ctx).Save(comment).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Comment{}, id).Error
}

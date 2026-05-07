package task

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uint) (*Task, error)
	GetByProjectID(ctx context.Context, pID uuid.UUID) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, task *Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *repository) GetByProjectID(ctx context.Context, pID uuid.UUID) ([]Task, error) {
	var tasks []Task
	err := r.db.WithContext(ctx).Where("project_id = ?", pID).Find(&tasks).Error
	return tasks, err
}

func (r *repository) Update(ctx context.Context, task *Task) error {
	return r.db.Save(task).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

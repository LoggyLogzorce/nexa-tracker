package task

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id uint, archived bool) (*Task, error)
	GetByProjectID(ctx context.Context, pID uuid.UUID, archived bool) ([]Task, error)
	Update(ctx context.Context, task *Task, history *UpdateHistory) error
	Delete(ctx context.Context, id uint) error

	GetHistoryByTaskID(ctx context.Context, taskID uint) ([]UpdateHistory, error)
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

func (r *repository) GetByID(ctx context.Context, id uint, archived bool) (*Task, error) {
	var task Task
	err := r.db.WithContext(ctx).Where("is_archive = ?", archived).First(&task, id).Error
	return &task, err
}

func (r *repository) GetByProjectID(ctx context.Context, pID uuid.UUID, archived bool) ([]Task, error) {
	var tasks []Task
	err := r.db.WithContext(ctx).Where("project_id = ? AND is_archive = ?", pID, archived).Find(&tasks).Error
	return tasks, err
}

func (r *repository) Update(ctx context.Context, task *Task, history *UpdateHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(task).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&Task{}, id).Error
}

func (r *repository) GetHistoryByTaskID(ctx context.Context, taskID uint) ([]UpdateHistory, error) {
	var history []UpdateHistory
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&history).Error
	return history, err
}

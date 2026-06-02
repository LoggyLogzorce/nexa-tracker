package task

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
	"time"
)

type Repository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uint, archived bool) (*models.Task, error)
	GetByProjectID(ctx context.Context, pID uuid.UUID, archived bool) ([]models.Task, error)
	GetByAssigneeID(ctx context.Context, uID uuid.UUID, archived bool) ([]models.Task, error)
	GetByReporterID(ctx context.Context, uID uuid.UUID, archived bool) ([]models.Task, error)
	GetByProjectIDAndUserID(ctx context.Context, pID uuid.UUID, userID uuid.UUID) ([]models.Task, error)
	Search(ctx context.Context, q string, projectIDs []uuid.UUID, limit int) ([]models.Task, error)
	Update(ctx context.Context, task *models.Task, history *models.UpdateHistory) error
	Delete(ctx context.Context, id uint) error
	DeleteParticipantInTask(ctx context.Context, tasks []models.Task, histories []models.UpdateHistory) error

	GetHistoryByTaskID(ctx context.Context, taskID uint) ([]models.UpdateHistory, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, task *models.Task) error {
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	return r.db.WithContext(ctx).Create(task).Error
}

func (r *repository) GetByID(ctx context.Context, id uint, archived bool) (*models.Task, error) {
	var task models.Task
	err := r.db.WithContext(ctx).Where("is_archive = ?", archived).First(&task, id).Error
	return &task, err
}

func (r *repository) GetByProjectID(ctx context.Context, pID uuid.UUID, archived bool) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).Where("project_id = ? AND is_archive = ?", pID, archived).Find(&tasks).Error
	return tasks, err
}

func (r *repository) GetByAssigneeID(ctx context.Context, uID uuid.UUID, archived bool) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("assignee_id = ? and is_archive = ?", uID, archived).
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) GetByReporterID(ctx context.Context, uID uuid.UUID, archived bool) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("reporter_id = ? and is_archive = ?", uID, archived).
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) GetByProjectIDAndUserID(ctx context.Context, pID uuid.UUID, userID uuid.UUID) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).Where("project_id = ? AND (assignee_id = ? OR reporter_id = ?)",
		pID, userID, userID).Find(&tasks).Error
	return tasks, err
}

func (r *repository) Search(ctx context.Context, q string, projectIDs []uuid.UUID, limit int) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.WithContext(ctx).
		Where("is_archive = ?", false).
		Where("project_id IN ?", projectIDs).
		Where("title ILIKE ?", "%"+q+"%").
		Limit(limit).
		Find(&tasks).Error
	return tasks, err
}

func (r *repository) Update(ctx context.Context, task *models.Task, history *models.UpdateHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(task).Error; err != nil {
			return err
		}
		return tx.Create(history).Error
	})
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Task{}, id).Error
}

func (r *repository) GetHistoryByTaskID(ctx context.Context, taskID uint) ([]models.UpdateHistory, error) {
	var history []models.UpdateHistory
	err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Find(&history).Error
	return history, err
}

func (r *repository) DeleteParticipantInTask(ctx context.Context, tasks []models.Task, histories []models.UpdateHistory) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&tasks).Error; err != nil {
			return err
		}
		return tx.Create(histories).Error
	})
}

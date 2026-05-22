package priority

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
)

type Repository interface {
	Create(ctx context.Context, priority *models.Priority) error
	CreateBatch(ctx context.Context, priorities []models.Priority) error
	GetByID(ctx context.Context, id uint) (*models.Priority, error)
	GetByTitle(ctx context.Context, title string, projectID uuid.UUID) (*models.Priority, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Priority, error)
	GetListByIDs(ctx context.Context, ids []uint) ([]models.Priority, error)
	GetListByProjectsIDs(ctx context.Context, ids []uuid.UUID) ([]models.Priority, error)
	Update(ctx context.Context, priority *models.Priority) error
	Delete(ctx context.Context, id uint) error
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, priority *models.Priority) error {
	return r.db.WithContext(ctx).Create(priority).Error
}

func (r *repository) CreateBatch(ctx context.Context, priorities []models.Priority) error {
	return r.db.WithContext(ctx).Create(&priorities).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Priority, error) {
	var priority models.Priority
	result := r.db.WithContext(ctx).First(&priority, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &priority, nil
}

func (r *repository) GetByTitle(ctx context.Context, title string, projectID uuid.UUID) (*models.Priority, error) {
	var priority models.Priority
	err := r.db.WithContext(ctx).Where("title = ? AND project_id = ?", title, projectID).First(&priority).Error
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Priority, error) {
	var priorities []models.Priority
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("id ASC").Find(&priorities)
	if result.Error != nil {
		return nil, result.Error
	}
	return priorities, nil
}

func (r *repository) GetListByIDs(ctx context.Context, ids []uint) ([]models.Priority, error) {
	var priorities []models.Priority
	err := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&priorities).Error
	return priorities, err
}

func (r *repository) GetListByProjectsIDs(ctx context.Context, ids []uuid.UUID) ([]models.Priority, error) {
	var priorities []models.Priority
	err := r.db.WithContext(ctx).Where("project_id IN (?)", ids).Find(&priorities).Error
	return priorities, err
}

func (r *repository) Update(ctx context.Context, priority *models.Priority) error {
	return r.db.WithContext(ctx).Save(priority).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Priority{}).Error
}

func (r *repository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.Priority{}).Error
}

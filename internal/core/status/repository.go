package status

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
)

type Repository interface {
	Create(ctx context.Context, status *models.Status) error
	CreateBatch(ctx context.Context, statuses []models.Status) error
	GetByID(ctx context.Context, id uint) (*models.Status, error)
	GetByName(ctx context.Context, name string) (*models.Status, error)
	GetByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) (*models.Status, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Status, error)
	GetListByIDs(ctx context.Context, ids []uint) ([]models.Status, error)
	GetListByProjectsIDs(ctx context.Context, ids []uuid.UUID) ([]models.Status, error)
	GetMaxOrderIndex(ctx context.Context, projectID uuid.UUID) (int, error)
	Update(ctx context.Context, status *models.Status) error
	UpdateOrderIndexBatch(ctx context.Context, updates []struct {
		ID         uint
		OrderIndex int
	}) error
	Delete(ctx context.Context, id uint) error
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, status *models.Status) error {
	return r.db.WithContext(ctx).Create(status).Error
}

func (r *repository) CreateBatch(ctx context.Context, statuses []models.Status) error {
	return r.db.Create(&statuses).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*models.Status, error) {
	var status models.Status
	result := r.db.WithContext(ctx).First(&status, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &status, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*models.Status, error) {
	var status *models.Status
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&status).Error
	return status, err
}

func (r *repository) GetByNameAndProjectID(ctx context.Context, name string, projectID uuid.UUID) (*models.Status, error) {
	var status *models.Status
	err := r.db.WithContext(ctx).Where("name = ? AND project_id = ?", name, projectID).First(&status).Error
	return status, err
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.Status, error) {
	var statuses []models.Status
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("order_index ASC").Find(&statuses)
	if result.Error != nil {
		return nil, result.Error
	}
	return statuses, nil
}

func (r *repository) GetListByIDs(ctx context.Context, ids []uint) ([]models.Status, error) {
	var statuses []models.Status
	err := r.db.WithContext(ctx).Where("id IN (?)", ids).Find(&statuses).Error
	return statuses, err
}

func (r *repository) GetListByProjectsIDs(ctx context.Context, ids []uuid.UUID) ([]models.Status, error) {
	var statuses []models.Status
	err := r.db.WithContext(ctx).Where("project_id IN (?)", ids).Find(&statuses).Error
	return statuses, err
}

func (r *repository) GetMaxOrderIndex(ctx context.Context, projectID uuid.UUID) (int, error) {
	var maxIndex int
	result := r.db.WithContext(ctx).Model(&models.Status{}).
		Where("project_id = ?", projectID).
		Select("COALESCE(MAX(order_index), -1)").
		Scan(&maxIndex)
	if result.Error != nil {
		return 0, result.Error
	}
	return maxIndex, nil
}

func (r *repository) Update(ctx context.Context, status *models.Status) error {
	return r.db.WithContext(ctx).Save(status).Error
}

func (r *repository) UpdateOrderIndexBatch(ctx context.Context, updates []struct {
	ID         uint
	OrderIndex int
}) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			if err := tx.Model(&models.Status{}).Where("id = ?", update.ID).Update("order_index", update.OrderIndex).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Status{}).Error
}

func (r *repository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&models.Status{}).Error
}

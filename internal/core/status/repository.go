package status

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, status *Status) error
	CreateBatch(ctx context.Context, statuses []Status) error
	GetByID(ctx context.Context, id uint) (*Status, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Status, error)
	Update(ctx context.Context, status *Status) error
	Delete(ctx context.Context, id uint) error
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, status *Status) error {
	// TODO: Implement
	return nil
}

func (r *repository) CreateBatch(ctx context.Context, statuses []Status) error {
	return r.db.Create(&statuses).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Status, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Status, error) {
	var statuses []Status
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("order_index ASC").Find(&statuses)
	if result.Error != nil {
		return nil, result.Error
	}
	return statuses, nil
}

func (r *repository) Update(ctx context.Context, status *Status) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	// TODO: Implement
	return nil
}

func (r *repository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&Status{}).Error
}

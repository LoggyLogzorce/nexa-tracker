package priority

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, priority *Priority) error
	CreateBatch(ctx context.Context, priorities []Priority) error
	GetByID(ctx context.Context, id uint) (*Priority, error)
	GetByTitle(ctx context.Context, title string, projectID uuid.UUID) (*Priority, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Priority, error)
	Update(ctx context.Context, priority *Priority) error
	Delete(ctx context.Context, id uint) error
	DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, priority *Priority) error {
	return r.db.WithContext(ctx).Create(priority).Error
}

func (r *repository) CreateBatch(ctx context.Context, priorities []Priority) error {
	return r.db.WithContext(ctx).Create(&priorities).Error
}

func (r *repository) GetByID(ctx context.Context, id uint) (*Priority, error) {
	var priority Priority
	result := r.db.WithContext(ctx).First(&priority, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &priority, nil
}

func (r *repository) GetByTitle(ctx context.Context, title string, projectID uuid.UUID) (*Priority, error) {
	var priority Priority
	err := r.db.WithContext(ctx).Where("title = ? AND project_id = ?", title, projectID).First(&priority).Error
	if err != nil {
		return nil, err
	}
	return &priority, nil
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]Priority, error) {
	var priorities []Priority
	result := r.db.WithContext(ctx).Where("project_id = ?", projectID).Order("id ASC").Find(&priorities)
	if result.Error != nil {
		return nil, result.Error
	}
	return priorities, nil
}

func (r *repository) Update(ctx context.Context, priority *Priority) error {
	return r.db.WithContext(ctx).Save(priority).Error
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&Priority{}).Error
}

func (r *repository) DeleteByProjectID(ctx context.Context, projectID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("project_id = ?", projectID).Delete(&Priority{}).Error
}

package project

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]Project, error)
	ListByParticipant(ctx context.Context, userID uuid.UUID, projectIDs []uuid.UUID) ([]Project, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, project *Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *repository) GetByID(ctx context.Context, id uuid.UUID) (*Project, error) {
	var project Project
	result := r.db.WithContext(ctx).Where("id = ?", id).First(&project)

	if result.Error != nil {
		return nil, result.Error
	}

	return &project, nil
}

func (r *repository) ListByOwner(ctx context.Context, ownerID uuid.UUID) ([]Project, error) {
	var projects []Project
	err := r.db.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&projects).Error
	return projects, err
}

func (r *repository) ListByParticipant(ctx context.Context, userID uuid.UUID, projectIDs []uuid.UUID) ([]Project, error) {
	var projects []Project
	err := r.db.WithContext(ctx).Where("owner_id = ? OR id IN ?", userID, projectIDs).Find(&projects).Error
	return projects, err
}

func (r *repository) Update(ctx context.Context, project *Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&Project{}, "id = ?", id).Error
}

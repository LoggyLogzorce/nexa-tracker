package project

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, project *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	List(ctx context.Context, ownerID uuid.UUID) ([]Project, error)
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

func (r *repository) List(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	var projects []Project
	result := r.db.WithContext(ctx).
		Joins("LEFT JOIN project_participants ON project_participants.project_id = projects.id").
		Where("project_participants.user_id = ? OR projects.owner_id = ?", userID, userID).
		Distinct().
		Find(&projects)

	if result.Error != nil {
		return nil, result.Error
	}

	return projects, nil
}

func (r *repository) Update(ctx context.Context, project *Project) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	// TODO: Implement
	return nil
}

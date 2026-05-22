package participant

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/models"
)

type Repository interface {
	Create(ctx context.Context, participant *models.ProjectParticipant) error
	GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.ProjectParticipant, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.ProjectParticipant, error)
	GetByProjectAndUser(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*models.ProjectParticipant, error)
	Update(ctx context.Context, participant *models.ProjectParticipant) error
	Delete(ctx context.Context, participant *models.ProjectParticipant) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, participant *models.ProjectParticipant) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *repository) GetByProjectID(ctx context.Context, projectID uuid.UUID) ([]models.ProjectParticipant, error) {
	var participants []models.ProjectParticipant
	err := r.db.WithContext(ctx).Where("project_id = ?", projectID).Find(&participants).Error
	return participants, err
}

func (r *repository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]models.ProjectParticipant, error) {
	var participants []models.ProjectParticipant
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&participants).Error
	return participants, err
}

func (r *repository) GetByProjectAndUser(ctx context.Context, projectID uuid.UUID, userID uuid.UUID) (*models.ProjectParticipant, error) {
	var participant models.ProjectParticipant

	result := r.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", projectID, userID).First(&participant)
	if result.Error != nil {
		return nil, result.Error
	}

	return &participant, nil
}

func (r *repository) Update(ctx context.Context, participant *models.ProjectParticipant) error {
	return r.db.WithContext(ctx).Save(participant).Error
}

func (r *repository) Delete(ctx context.Context, participant *models.ProjectParticipant) error {
	return r.db.WithContext(ctx).Where("project_id = ? AND user_id = ?", participant.ProjectID, participant.UserID).
		Delete(&models.ProjectParticipant{}).Error
}

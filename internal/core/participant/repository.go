package participant

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(participant *ProjectParticipant) error
	GetByProjectID(projectID uuid.UUID) ([]ProjectParticipant, error)
	GetByUserID(userID string) ([]ProjectParticipant, error)
	GetByProjectAndUser(projectID uuid.UUID, userID string) (*ProjectParticipant, error)
	Update(participant *ProjectParticipant) error
	Delete(projectID uuid.UUID, userID string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(participant *ProjectParticipant) error {
	// TODO: Implement
	return nil
}

func (r *repository) GetByProjectID(projectID uuid.UUID) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByUserID(userID string) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectAndUser(projectID uuid.UUID, userID string) (*ProjectParticipant, error) {
	var participant ProjectParticipant

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	result := r.db.Where("project_id = ? AND user_id = ?", projectID, userUUID).First(&participant)
	if result.Error != nil {
		return nil, result.Error
	}

	return &participant, nil
}

func (r *repository) Update(participant *ProjectParticipant) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(projectID uuid.UUID, userID string) error {
	// TODO: Implement
	return nil
}

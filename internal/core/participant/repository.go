package participant

import "gorm.io/gorm"

type Repository interface {
	Create(participant *ProjectParticipant) error
	GetByProjectID(projectID uint) ([]ProjectParticipant, error)
	GetByUserID(userID string) ([]ProjectParticipant, error)
	GetByProjectAndUser(projectID uint, userID string) (*ProjectParticipant, error)
	Update(participant *ProjectParticipant) error
	Delete(projectID uint, userID string) error
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

func (r *repository) GetByProjectID(projectID uint) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByUserID(userID string) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) GetByProjectAndUser(projectID uint, userID string) (*ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (r *repository) Update(participant *ProjectParticipant) error {
	// TODO: Implement
	return nil
}

func (r *repository) Delete(projectID uint, userID string) error {
	// TODO: Implement
	return nil
}

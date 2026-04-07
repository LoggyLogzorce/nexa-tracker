package participant

type Service interface {
	AddParticipant(participant *ProjectParticipant) error
	GetByProjectID(projectID uint) ([]ProjectParticipant, error)
	GetByUserID(userID string) ([]ProjectParticipant, error)
	UpdateRole(projectID uint, userID string, role string) error
	RemoveParticipant(projectID uint, userID string) error
	CheckAccess(projectID uint, userID string, requiredRole string) (bool, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) AddParticipant(participant *ProjectParticipant) error {
	// TODO: Implement
	return nil
}

func (s *service) GetByProjectID(projectID uint) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByUserID(userID string) ([]ProjectParticipant, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) UpdateRole(projectID uint, userID string, role string) error {
	// TODO: Implement
	return nil
}

func (s *service) RemoveParticipant(projectID uint, userID string) error {
	// TODO: Implement
	return nil
}

func (s *service) CheckAccess(projectID uint, userID string, requiredRole string) (bool, error) {
	// TODO: Implement role hierarchy check
	return false, nil
}

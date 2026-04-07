package history

type Service interface {
	GetByTaskID(taskID uint) ([]UpdateHistory, error)
	GetByUserID(userID string) ([]UpdateHistory, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByTaskID(taskID uint) ([]UpdateHistory, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByUserID(userID string) ([]UpdateHistory, error) {
	// TODO: Implement
	return nil, nil
}

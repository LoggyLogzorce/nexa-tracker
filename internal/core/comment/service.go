package comment

type Service interface {
	Create(comment *Comment) error
	GetByID(id uint) (*Comment, error)
	GetByTaskID(taskID uint) ([]Comment, error)
	Update(comment *Comment) error
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(comment *Comment) error {
	// TODO: Implement
	return nil
}

func (s *service) GetByID(id uint) (*Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByTaskID(taskID uint) ([]Comment, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Update(comment *Comment) error {
	// TODO: Implement
	return nil
}

func (s *service) Delete(id uint) error {
	// TODO: Implement
	return nil
}

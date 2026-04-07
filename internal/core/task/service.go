package task

type Service interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	List(filters map[string]interface{}) ([]Task, error)
	Update(task *Task) error
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(task *Task) error {
	return s.repo.Create(task)
}

func (s *service) GetByID(id uint) (*Task, error) {
	return s.repo.GetByID(id)
}

func (s *service) List(filters map[string]interface{}) ([]Task, error) {
	return s.repo.List(filters)
}

func (s *service) Update(task *Task) error {
	return s.repo.Update(task)
}

func (s *service) Delete(id uint) error {
	return s.repo.Delete(id)
}

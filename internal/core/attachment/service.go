package attachment

type Service interface {
	Upload(attachment *Attachment, fileData []byte) error
	GetByID(id uint) (*Attachment, error)
	GetByTaskID(taskID uint) ([]Attachment, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Upload(attachment *Attachment, fileData []byte) error {
	// TODO: Implement file upload to /var/data/nexa-tracker
	return nil
}

func (s *service) GetByID(id uint) (*Attachment, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) GetByTaskID(taskID uint) ([]Attachment, error) {
	// TODO: Implement
	return nil, nil
}

func (s *service) Delete(id uint) error {
	// TODO: Implement (delete file from disk + DB record)
	return nil
}

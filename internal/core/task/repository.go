package task

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(task *Task) error
	GetByID(id uint) (*Task, error)
	List(filters map[string]interface{}) ([]Task, error)
	Update(task *Task) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(task *Task) error {
	return r.db.Create(task).Error
}

func (r *repository) GetByID(id uint) (*Task, error) {
	var task Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *repository) List(filters map[string]interface{}) ([]Task, error) {
	var tasks []Task
	query := r.db.Model(&Task{})

	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}

	err := query.Find(&tasks).Error
	return tasks, err
}

func (r *repository) Update(task *Task) error {
	return r.db.Save(task).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Task{}, id).Error
}

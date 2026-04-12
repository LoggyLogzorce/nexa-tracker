package project

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, project *Project, ownerID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Project, error)
	List(ctx context.Context, userID uuid.UUID) ([]Project, error)
	Update(ctx context.Context, project *Project, userID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

type service struct {
	repo            Repository
	participantRepo participant.Repository
	eventBus        *events.EventBus
}

func NewService(repo Repository, eventBus *events.EventBus, participantRepo participant.Repository) Service {
	return &service{
		repo:            repo,
		participantRepo: participantRepo,
		eventBus:        eventBus,
	}
}

func (s *service) Create(ctx context.Context, project *Project, ownerID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Установить id и owner
	project.ID = uuid.New()
	project.OwnerID = ownerID

	// 2. Создать проект
	if err := s.repo.Create(ctxT, project); err != nil {
		return err
	}

	event := ProjectEvent{
		Type:      events.ProjectCreated,
		ProjectID: project.ID,
	}
	s.eventBus.Publish(event.ToEvent())

	return nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Project, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить проект из БД
	project, err := s.repo.GetByID(ctxT, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	//// 2. Проверить права доступа
	//if project.OwnerID == userID {
	//	// Пользователь - owner, доступ разрешен
	//	return project, nil
	//}
	//
	//// Проверить, является ли пользователь участником проекта
	//_, err = s.participantRepo.GetByProjectAndUser(id, userID.String())
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		// Пользователь не является ни owner'ом, ни участником
	//		return nil, ErrProjectAccessDenied
	//	}
	//	// Другая ошибка БД
	//	return nil, err
	//}

	// Пользователь является участником, доступ разрешен
	return project, nil
}

func (s *service) List(ctx context.Context, userID uuid.UUID) ([]Project, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	projects, err := s.repo.List(ctxT, userID)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *service) Update(ctx context.Context, project *Project, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 1. Получить проект из БД
	existingProject, err := s.repo.GetByID(ctxT, project.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	//// 2. Проверить права доступа - только owner может обновлять проект
	//if existingProject.OwnerID != userID {
	//	return ErrProjectAccessDenied
	//}

	if project.Title == "" {
		project.Title = existingProject.Title
	}

	if project.Description == nil {
		project.Description = existingProject.Description
	}

	project.ID = existingProject.ID
	project.OwnerID = existingProject.OwnerID
	project.CreatedAt = existingProject.CreatedAt

	// 3. Обновить проект
	if err := s.repo.Update(ctxT, existingProject); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}

	return nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//// 1. Получить проект из БД
	//existingProject, err := s.repo.GetByID(ctxT, id)
	//if err != nil {
	//	if errors.Is(err, gorm.ErrRecordNotFound) {
	//		return ErrProjectNotFound
	//	}
	//	return err
	//}
	//
	//// 2. Проверить права доступа - только owner может удалять проект
	//if existingProject.OwnerID != userID {
	//	return ErrProjectAccessDenied
	//}

	// 3. Удалить проект
	if err := s.repo.Delete(ctxT, id); err != nil {
		return err
	}

	// 4. Опубликовать событие ProjectDeleted
	event := ProjectEvent{
		Type:      events.ProjectDeleted,
		ProjectID: id,
	}
	s.eventBus.Publish(event.ToEvent())

	return nil
}

package task

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, task *Task) (*TaskResponse, error)
	GetByID(ctx context.Context, id uint, param string) (*TaskResponse, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, param string) ([]TaskResponse, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id uint) error
}

type service struct {
	repo            Repository
	userRepo        user.Repository
	statusRepo      status.Repository
	priorityRepo    priority.Repository
	participantRepo participant.Repository
	eventBus        events.EventBus
}

func NewService(repo Repository, userRepo user.Repository, statusRepo status.Repository, priorityRepo priority.Repository, participantRepo participant.Repository) Service {
	return &service{
		repo:            repo,
		userRepo:        userRepo,
		statusRepo:      statusRepo,
		priorityRepo:    priorityRepo,
		participantRepo: participantRepo,
	}
}

func (s *service) Create(ctx context.Context, task *Task) (*TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if task.AssigneeID != nil {
		assignee, err := s.participantRepo.GetByUserID(ctxT, *task.AssigneeID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDataIntegrity
		}

		if assignee == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAssigneeNotInProject
		}
	}

	taskRes := &TaskResponse{
		IsArchive: false,
	}

	var st *status.Status
	var err error
	if task.StatusID != nil {
		st, err = s.statusRepo.GetByID(ctxT, *task.StatusID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStatusNotInProject
		}
		if err != nil {
			return nil, err
		}
		if st.ProjectID != task.ProjectID {
			return nil, ErrStatusNotInProject
		}
		taskRes.Status = &TaskStatusResponse{
			ID:         st.ID,
			Name:       st.Name,
			Color:      st.Color,
			OrderIndex: st.OrderIndex,
		}
	}

	var pr *priority.Priority
	if task.PriorityID != nil {
		pr, err = s.priorityRepo.GetByID(ctxT, *task.PriorityID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPriorityNotInProject
		}
		if err != nil {
			return nil, err
		}
		if pr.ProjectID != task.ProjectID {
			return nil, ErrPriorityNotInProject
		}
		taskRes.Priority = &TaskPriorityResponse{
			ID:    pr.ID,
			Title: pr.Title,
			Color: pr.Color,
		}
	}

	userIDs := make([]uuid.UUID, 0, 2)
	if task.AssigneeID != nil {
		userIDs = append(userIDs, *task.AssigneeID)
	}
	if task.ReporterID != nil {
		userIDs = append(userIDs, *task.ReporterID)
	}

	if len(userIDs) > 0 {
		users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			if task.AssigneeID != nil && *task.AssigneeID == u.ID {
				taskRes.Assignee = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
			}
			if task.ReporterID != nil && *task.ReporterID == u.ID {
				taskRes.Reporter = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email}
			}
		}
	}

	err = s.repo.Create(ctx, task)
	if err != nil {
		return nil, err
	}

	if task.Deadline != nil {
		formatted := task.Deadline.Format("2006-01-02")
		taskRes.Deadline = &formatted
	}

	taskRes.ID = task.ID
	taskRes.Title = task.Title
	taskRes.ProjectID = task.ProjectID
	taskRes.Description = task.Description
	taskRes.CreatedAt = task.CreatedAt
	taskRes.UpdatedAt = task.UpdatedAt

	go func() {
		event := events.TaskEvent{
			Type:        events.TaskCreate,
			ID:          task.ID,
			CreatedAt:   task.CreatedAt,
			UpdatedAt:   task.UpdatedAt,
			Title:       task.Title,
			Description: task.Description,
			Deadline:    task.Deadline,
			StatusID:    task.StatusID,
			PriorityID:  task.PriorityID,
			AssigneeID:  task.AssigneeID,
			ReporterID:  task.ReporterID,
		}
		s.eventBus.Publish(event.ToEvent())
	}()

	return taskRes, nil
}

func (s *service) GetByID(ctx context.Context, id uint, param string) (*TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	archived := false
	if param == "true" {
		archived = true
	}

	task, err := s.repo.GetByID(ctxT, id, archived)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	taskRes := &TaskResponse{
		ID:          task.ID,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Title:       task.Title,
		Description: task.Description,
		ProjectID:   task.ProjectID,
	}

	if task.Deadline != nil {
		deadLine := task.Deadline.Format("2006-01-02")
		taskRes.Deadline = &deadLine
	}

	var status *status.Status
	if task.StatusID != nil {
		status, err = s.statusRepo.GetByID(ctxT, *task.StatusID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrStatusNotInProject
		}
		if err != nil {
			return nil, err
		}
		taskRes.Status = &TaskStatusResponse{
			ID:         status.ID,
			Name:       status.Name,
			Color:      status.Color,
			OrderIndex: status.OrderIndex,
		}
	}

	var priority *priority.Priority
	if task.PriorityID != nil {
		priority, err = s.priorityRepo.GetByID(ctxT, *task.PriorityID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPriorityNotInProject
		}
		if err != nil {
			return nil, err
		}
		taskRes.Priority = &TaskPriorityResponse{
			ID:    priority.ID,
			Title: priority.Title,
			Color: priority.Color,
		}
	}

	// Загружаем пользователей
	userIDs := make([]uuid.UUID, 0, 2)
	if task.AssigneeID != nil {
		userIDs = append(userIDs, *task.AssigneeID)
	}
	if task.ReporterID != nil {
		userIDs = append(userIDs, *task.ReporterID)
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]user.User, 2)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	if task.AssigneeID != nil {
		taskRes.Assignee = &TaskUserResponse{
			ID:    usersMap[*task.AssigneeID].ID,
			Name:  usersMap[*task.AssigneeID].Name,
			Email: usersMap[*task.AssigneeID].Email,
		}
	}

	if task.ReporterID != nil {
		taskRes.Reporter = &TaskUserResponse{
			ID:    usersMap[*task.ReporterID].ID,
			Name:  usersMap[*task.ReporterID].Name,
			Email: usersMap[*task.ReporterID].Email,
		}
	}

	return taskRes, nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID, param string) ([]TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	archived := false
	if param == "true" {
		archived = true
	}

	tasks, err := s.repo.GetByProjectID(ctxT, projectID, archived)
	if err != nil {
		return nil, err
	}

	// Собираем уникальные ID
	userIDsMap := make(map[uuid.UUID]struct{})
	statusIDsMap := make(map[uint]struct{})
	priorityIDsMap := make(map[uint]struct{})

	for _, t := range tasks {
		if t.AssigneeID != nil {
			userIDsMap[*t.AssigneeID] = struct{}{}
		}
		if t.ReporterID != nil {
			userIDsMap[*t.ReporterID] = struct{}{}
		}
		if t.StatusID != nil {
			statusIDsMap[*t.StatusID] = struct{}{}
		}
		if t.PriorityID != nil {
			priorityIDsMap[*t.PriorityID] = struct{}{}
		}
	}

	// Загружаем пользователей
	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]user.User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	// Загружаем статусы
	statusIDs := make([]uint, 0, len(statusIDsMap))
	for id := range statusIDsMap {
		statusIDs = append(statusIDs, id)
	}
	statuses, err := s.statusRepo.GetListByIDs(ctxT, statusIDs)
	if err != nil {
		return nil, err
	}
	statusesMap := make(map[uint]status.Status, len(statuses))
	for _, st := range statuses {
		statusesMap[st.ID] = st
	}

	// Загружаем приоритеты
	priorityIDs := make([]uint, 0, len(priorityIDsMap))
	for id := range priorityIDsMap {
		priorityIDs = append(priorityIDs, id)
	}
	priorities, err := s.priorityRepo.GetListByIDs(ctxT, priorityIDs)
	if err != nil {
		return nil, err
	}
	prioritiesMap := make(map[uint]priority.Priority, len(priorities))
	for _, p := range priorities {
		prioritiesMap[p.ID] = p
	}

	// Собираем ответ
	response := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = TaskResponse{
			ID:          t.ID,
			CreatedAt:   t.CreatedAt,
			UpdatedAt:   t.UpdatedAt,
			Title:       t.Title,
			Description: t.Description,
			ProjectID:   t.ProjectID,
			IsArchive:   t.IsArchive,
		}
		if t.Deadline != nil {
			formatted := t.Deadline.Format("2006-01-02")
			response[i].Deadline = &formatted
		}

		if t.AssigneeID != nil {
			if u, ok := usersMap[*t.AssigneeID]; ok {
				response[i].Assignee = &TaskUserResponse{ID: u.ID,
					Name:  u.Name,
					Email: u.Email,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}

		if t.ReporterID != nil {
			if u, ok := usersMap[*t.ReporterID]; ok {
				response[i].Reporter = &TaskUserResponse{
					ID:    u.ID,
					Name:  u.Name,
					Email: u.Email,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}

		if t.StatusID != nil {
			if st, ok := statusesMap[*t.StatusID]; ok {
				response[i].Status = &TaskStatusResponse{
					ID:         st.ID,
					Name:       st.Name,
					Color:      st.Color,
					OrderIndex: st.OrderIndex,
				}
			} else {
				// TODO усправить
				return nil, ErrDataIntegrity
			}
		}

		if t.PriorityID != nil {
			if p, ok := prioritiesMap[*t.PriorityID]; ok {
				response[i].Priority = &TaskPriorityResponse{
					ID:    p.ID,
					Title: p.Title,
					Color: p.Color,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}
	}

	return response, nil
}

func (s *service) Update(ctx context.Context, task *Task) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Update(ctxT, task)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return s.repo.Delete(ctxT, id)
}

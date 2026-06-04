package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"log"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/models"
	events2 "nexa-task-tracker/pkg/events"
	"time"
)

type Service interface {
	Create(ctx context.Context, task *models.Task) (*TaskResponse, error)
	GetByID(ctx context.Context, id uint, param string) (*TaskResponse, error)
	GetByProjectID(ctx context.Context, projectID uuid.UUID, param Param) ([]TaskResponse, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, param Param) ([]TaskResponse, error)
	Search(ctx context.Context, q string, userID uuid.UUID) ([]TaskResponse, error)
	Update(ctx context.Context, taskID uint, req *UpdateTaskRequest, archived bool, userID uuid.UUID) (*TaskResponse, error)
	Delete(ctx context.Context, taskId uint, userID uuid.UUID) error

	GetHistoryByTaskID(ctx context.Context, taskID uint) ([]HistoryResponse, error)

	HandleParticipantDelete(event events2.Event) error
}

type service struct {
	repo            Repository
	userRepo        user.Repository
	statusRepo      status.Repository
	priorityRepo    priority.Repository
	participantRepo participant.Repository
	projectRepo     project.Repository
	eventBus        *events2.EventBus
}

func NewService(repo Repository, userRepo user.Repository, statusRepo status.Repository, priorityRepo priority.Repository, participantRepo participant.Repository, projectRepo project.Repository, eventBus *events2.EventBus) Service {
	return &service{
		repo:            repo,
		userRepo:        userRepo,
		statusRepo:      statusRepo,
		priorityRepo:    priorityRepo,
		participantRepo: participantRepo,
		projectRepo:     projectRepo,
		eventBus:        eventBus,
	}
}

type FieldChange struct {
	Field    string `json:"field"`
	OldValue any    `json:"old_value"`
	NewValue any    `json:"new_value"`
	OldName  string `json:"old_name,omitempty"`
	NewName  string `json:"new_name,omitempty"`
}

func (s *service) Create(ctx context.Context, task *models.Task) (*TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if task.AssigneeID != nil {
		assignee, err := s.participantRepo.GetByProjectAndUser(ctxT, task.ProjectID, *task.AssigneeID)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDataIntegrity
		}

		if assignee == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			projectData, err := s.projectRepo.GetByID(ctxT, task.ProjectID)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrProjectNotFound
			}
			if projectData.OwnerID != *task.AssigneeID {
				return nil, ErrAssigneeNotInProject
			}
		} else if assignee.Role == "read_only" {
			return nil, ErrInvalidAssigneeRole
		}
	}

	taskRes := &TaskResponse{
		IsArchive: false,
	}

	var st *models.Status
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

	var pr *models.Priority
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
				taskRes.Assignee = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email, AvatarUrl: u.AvatarUrl}
			}
			if task.ReporterID != nil && *task.ReporterID == u.ID {
				taskRes.Reporter = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email, AvatarUrl: u.AvatarUrl}
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
		event := events2.TaskEvent{
			Type:        events2.TaskCreate,
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
		IsArchive:   task.IsArchive,
	}

	if task.Deadline != nil {
		deadLine := task.Deadline.Format("2006-01-02")
		taskRes.Deadline = &deadLine
	}

	var status *models.Status
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

	var priority *models.Priority
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
	usersMap := make(map[uuid.UUID]models.User, 2)
	for _, u := range users {
		usersMap[u.ID] = u
	}

	if task.AssigneeID != nil {
		taskRes.Assignee = &TaskUserResponse{
			ID:        usersMap[*task.AssigneeID].ID,
			Name:      usersMap[*task.AssigneeID].Name,
			Email:     usersMap[*task.AssigneeID].Email,
			AvatarUrl: usersMap[*task.AssigneeID].AvatarUrl,
		}
	}

	if task.ReporterID != nil {
		taskRes.Reporter = &TaskUserResponse{
			ID:        usersMap[*task.ReporterID].ID,
			Name:      usersMap[*task.ReporterID].Name,
			Email:     usersMap[*task.ReporterID].Email,
			AvatarUrl: usersMap[*task.ReporterID].AvatarUrl,
		}
	}

	return taskRes, nil
}

func (s *service) GetByProjectID(ctx context.Context, projectID uuid.UUID, param Param) ([]TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	archived := false
	if param.Archived == "true" {
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
	usersMap := make(map[uuid.UUID]models.User, len(users))
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
	statusesMap := make(map[uint]models.Status, len(statuses))
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
	prioritiesMap := make(map[uint]models.Priority, len(priorities))
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
				response[i].Assignee = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}

		if t.ReporterID != nil {
			if u, ok := usersMap[*t.ReporterID]; ok {
				response[i].Reporter = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
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

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID, param Param) ([]TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var tasks []models.Task
	var err error
	archived := false
	if param.Archived == "true" {
		archived = true
	}
	if param.UserParam == "assigned" {
		tasks, err = s.repo.GetByAssigneeID(ctxT, userID, archived)
		if err != nil {
			return nil, err
		}
	} else if param.UserParam == "reported" {
		tasks, err = s.repo.GetByReporterID(ctxT, userID, archived)
		if err != nil {
			return nil, err
		}
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
	usersMap := make(map[uuid.UUID]models.User, len(users))
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
	statusesMap := make(map[uint]models.Status, len(statuses))
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
	prioritiesMap := make(map[uint]models.Priority, len(priorities))
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
				response[i].Assignee = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}

		if t.ReporterID != nil {
			if u, ok := usersMap[*t.ReporterID]; ok {
				response[i].Reporter = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
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

func (s *service) Search(ctx context.Context, q string, userID uuid.UUID) ([]TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if len(q) < 2 {
		return []TaskResponse{}, nil
	}

	participants, err := s.participantRepo.GetByUserID(ctxT, userID)
	if err != nil {
		return nil, err
	}

	projectIDs := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		projectIDs[i] = p.ProjectID
	}

	ownedProjects, err := s.projectRepo.ListByOwner(ctxT, userID)
	if err != nil {
		return nil, err
	}
	ownedMap := make(map[uuid.UUID]bool)
	for _, p := range ownedProjects {
		ownedMap[p.ID] = true
		found := false
		for _, id := range projectIDs {
			if id == p.ID {
				found = true
				break
			}
		}
		if !found {
			projectIDs = append(projectIDs, p.ID)
		}
	}

	if len(projectIDs) == 0 {
		return []TaskResponse{}, nil
	}

	tasks, err := s.repo.Search(ctxT, q, projectIDs, 10)
	if err != nil {
		return nil, err
	}

	projectTitleMap := make(map[uuid.UUID]string, len(ownedProjects))
	for _, p := range ownedProjects {
		projectTitleMap[p.ID] = p.Title
	}
	for _, p := range participants {
		projectTitleMap[p.ProjectID] = ""
	}

	projectIDsForFetch := make([]uuid.UUID, 0)
	for _, t := range tasks {
		if _, ok := projectTitleMap[t.ProjectID]; !ok {
			projectIDsForFetch = append(projectIDsForFetch, t.ProjectID)
		}
	}
	if len(projectIDsForFetch) > 0 {
		for _, pID := range projectIDsForFetch {
			p, err := s.projectRepo.GetByID(ctxT, pID)
			if err == nil {
				projectTitleMap[pID] = p.Title
			}
		}
	}

	userIDsMap := make(map[uuid.UUID]struct{})
	for _, t := range tasks {
		if t.AssigneeID != nil {
			userIDsMap[*t.AssigneeID] = struct{}{}
		}
		if t.ReporterID != nil {
			userIDsMap[*t.ReporterID] = struct{}{}
		}
	}

	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		userIDs = append(userIDs, id)
	}
	usersMap := make(map[uuid.UUID]models.User)
	if len(userIDs) > 0 {
		users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
		if err == nil {
			for _, u := range users {
				usersMap[u.ID] = u
			}
		}
	}

	response := make([]TaskResponse, len(tasks))
	for i, t := range tasks {
		response[i] = TaskResponse{
			ID:           t.ID,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
			Title:        t.Title,
			Description:  t.Description,
			ProjectID:    t.ProjectID,
			ProjectTitle: projectTitleMap[t.ProjectID],
			IsArchive:    t.IsArchive,
		}
		if t.Deadline != nil {
			formatted := t.Deadline.Format("2006-01-02")
			response[i].Deadline = &formatted
		}
		if t.AssigneeID != nil {
			if u, ok := usersMap[*t.AssigneeID]; ok {
				response[i].Assignee = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
				}
			}
		}
		if t.ReporterID != nil {
			if u, ok := usersMap[*t.ReporterID]; ok {
				response[i].Reporter = &TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
				}
			}
		}
	}

	return response, nil
}

func (s *service) Update(ctx context.Context, taskID uint, req *UpdateTaskRequest, archived bool, userID uuid.UUID) (*TaskResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 7*time.Second)
	defer cancel()

	taskOld, err := s.repo.GetByID(ctxT, taskID, archived)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	taskNew := *taskOld
	var changes []FieldChange

	// title
	if req.Title != nil && *req.Title != taskOld.Title {
		changes = append(changes, FieldChange{Field: "title", OldValue: taskOld.Title, NewValue: *req.Title})
		taskNew.Title = *req.Title
	}

	// description
	if req.Description.Set && !equalStringPtr(req.Description.Value, taskOld.Description) {
		changes = append(changes, FieldChange{Field: "description", OldValue: taskOld.Description, NewValue: req.Description.Value})
		taskNew.Description = req.Description.Value
	}

	// status
	if req.StatusID.Set && !equalUintPtr(req.StatusID.Value, taskOld.StatusID) {
		if req.StatusID.Value != nil {
			st, err := s.statusRepo.GetByID(ctxT, *req.StatusID.Value)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrStatusNotInProject
				}
				return nil, err
			}
			if st.ProjectID != taskNew.ProjectID {
				return nil, ErrStatusNotInProject
			}
		}
		fc := FieldChange{Field: "status", OldValue: taskOld.StatusID, NewValue: req.StatusID.Value}
		if taskOld.StatusID != nil {
			if oldSt, err := s.statusRepo.GetByID(ctxT, *taskOld.StatusID); err == nil {
				fc.OldName = oldSt.Name
			}
		}
		if req.StatusID.Value != nil {
			if newSt, err := s.statusRepo.GetByID(ctxT, *req.StatusID.Value); err == nil {
				fc.NewName = newSt.Name
			}
		}
		changes = append(changes, fc)
		taskNew.StatusID = req.StatusID.Value
	}

	// priority
	if req.PriorityID.Set && !equalUintPtr(req.PriorityID.Value, taskOld.PriorityID) {
		if req.PriorityID.Value != nil {
			pr, err := s.priorityRepo.GetByID(ctxT, *req.PriorityID.Value)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, ErrPriorityNotInProject
				}
				return nil, err
			}
			if pr.ProjectID != taskNew.ProjectID {
				return nil, ErrPriorityNotInProject
			}
		}
		fc := FieldChange{Field: "priority", OldValue: taskOld.PriorityID, NewValue: req.PriorityID.Value}
		if taskOld.PriorityID != nil {
			if oldPr, err := s.priorityRepo.GetByID(ctxT, *taskOld.PriorityID); err == nil {
				fc.OldName = oldPr.Title
			}
		}
		if req.PriorityID.Value != nil {
			if newPr, err := s.priorityRepo.GetByID(ctxT, *req.PriorityID.Value); err == nil {
				fc.NewName = newPr.Title
			}
		}
		changes = append(changes, fc)
		taskNew.PriorityID = req.PriorityID.Value
	}

	// assignee
	if req.AssigneeID.Set && !equalUUIDPtr(req.AssigneeID.Value, taskOld.AssigneeID) {
		if req.AssigneeID.Value != nil {
			assignee, err := s.participantRepo.GetByProjectAndUser(ctxT, taskNew.ProjectID, *req.AssigneeID.Value)
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			if assignee == nil || errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrAssigneeNotInProject
			}
			if assignee.Role != "member" {
				return nil, ErrInvalidAssigneeRole
			}
		}
		fc := FieldChange{Field: "assignee", OldValue: taskOld.AssigneeID, NewValue: req.AssigneeID.Value}
		if taskOld.AssigneeID != nil {
			if oldUser, err := s.userRepo.GetByID(ctxT, *taskOld.AssigneeID); err == nil {
				fc.OldName = oldUser.Name
			}
		}
		if req.AssigneeID.Value != nil {
			if newUser, err := s.userRepo.GetByID(ctxT, *req.AssigneeID.Value); err == nil {
				fc.NewName = newUser.Name
			}
		}
		changes = append(changes, fc)
		taskNew.AssigneeID = req.AssigneeID.Value
	}

	// deadline
	if req.Deadline.Set && !equalTimePtr(req.Deadline.Value, taskOld.Deadline) {
		changes = append(changes, FieldChange{Field: "deadline", OldValue: taskOld.Deadline, NewValue: req.Deadline.Value})
		taskNew.Deadline = req.Deadline.Value
	}

	// is_archive
	if req.IsArchive.Set && req.IsArchive.Value != nil && *req.IsArchive.Value != taskOld.IsArchive {
		changes = append(changes, FieldChange{Field: "is_archive", OldValue: taskOld.IsArchive, NewValue: *req.IsArchive.Value})
		taskNew.IsArchive = *req.IsArchive.Value
	}

	if len(changes) == 0 {
		return nil, ErrNoFieldsToUpdate
	}

	// сборка taskRes
	taskRes := &TaskResponse{
		ID:          taskNew.ID,
		CreatedAt:   taskNew.CreatedAt,
		Title:       taskNew.Title,
		Description: taskNew.Description,
		ProjectID:   taskNew.ProjectID,
		IsArchive:   taskNew.IsArchive,
	}
	if taskNew.Deadline != nil {
		deadline := taskNew.Deadline.Format("2006-01-02")
		taskRes.Deadline = &deadline
	}
	if taskNew.StatusID != nil {
		if st, err := s.statusRepo.GetByID(ctxT, *taskNew.StatusID); err == nil {
			taskRes.Status = &TaskStatusResponse{
				ID:         st.ID,
				Name:       st.Name,
				Color:      st.Color,
				OrderIndex: st.OrderIndex,
			}
		}
	}
	if taskNew.PriorityID != nil {
		if pr, err := s.priorityRepo.GetByID(ctxT, *taskNew.PriorityID); err == nil {
			taskRes.Priority = &TaskPriorityResponse{
				ID:    pr.ID,
				Title: pr.Title,
				Color: pr.Color,
			}
		}
	}

	userIDs := make([]uuid.UUID, 0, 2)
	if taskNew.AssigneeID != nil {
		userIDs = append(userIDs, *taskNew.AssigneeID)
	}
	if taskNew.ReporterID != nil {
		userIDs = append(userIDs, *taskNew.ReporterID)
	}
	if len(userIDs) > 0 {
		users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
		if err != nil {
			return nil, err
		}
		for _, u := range users {
			if taskNew.AssigneeID != nil && *taskNew.AssigneeID == u.ID {
				taskRes.Assignee = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email, AvatarUrl: u.AvatarUrl}
			}
			if taskNew.ReporterID != nil && *taskNew.ReporterID == u.ID {
				taskRes.Reporter = &TaskUserResponse{ID: u.ID, Name: u.Name, Email: u.Email, AvatarUrl: u.AvatarUrl}
			}
		}
	}

	// history
	oldJSON, err := json.Marshal(taskOld)
	if err != nil {
		return nil, err
	}
	newJSON, err := json.Marshal(taskNew)
	if err != nil {
		return nil, err
	}
	changesJSON, err := json.Marshal(changes)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	history := &models.UpdateHistory{
		CreatedAt: now,
		UserID:    &userID,
		TaskID:    taskNew.ID,
		Old:       datatypes.JSON(oldJSON),
		New:       datatypes.JSON(newJSON),
		Changes:   datatypes.JSON(changesJSON),
	}
	taskNew.UpdatedAt = now
	taskRes.UpdatedAt = now

	if err := s.repo.Update(ctxT, &taskNew, history); err != nil {
		return nil, err
	}

	return taskRes, nil
}

func (s *service) Delete(ctx context.Context, taskId uint, userID uuid.UUID) error {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := s.repo.Delete(ctxT, taskId)
	if err != nil {
		return err
	}

	go func() {
		event := events2.TaskEvent{
			Type:      events2.TaskDelete,
			ID:        taskId,
			DeletedBy: userID,
		}
		s.eventBus.Publish(event.ToEvent())
	}()

	return nil
}

func (s *service) GetHistoryByTaskID(ctx context.Context, taskID uint) ([]HistoryResponse, error) {
	ctxT, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	history, err := s.repo.GetHistoryByTaskID(ctxT, taskID)
	if err != nil {
		return nil, err
	}

	userIDsMap := make(map[*uuid.UUID]struct{})

	for _, h := range history {
		userIDsMap[h.UserID] = struct{}{}
	}

	// Загружаем пользователей
	userIDs := make([]uuid.UUID, 0, len(userIDsMap))
	for id := range userIDsMap {
		if id != nil {
			userIDs = append(userIDs, *id)
		}
	}
	users, err := s.userRepo.GetListByIDs(ctxT, userIDs)
	if err != nil {
		return nil, err
	}
	usersMap := make(map[uuid.UUID]models.User, len(users))
	for _, u := range users {
		usersMap[u.ID] = u
	}

	response := make([]HistoryResponse, len(history))
	for i, h := range history {
		var changes []FieldChange
		if err := json.Unmarshal(h.Changes, &changes); err != nil {
			return nil, err
		}
		response[i] = HistoryResponse{
			ID:        h.ID,
			CreatedAt: h.CreatedAt,
			TaskID:    h.TaskID,
			Old:       h.Old,
			New:       h.New,
			Changes:   changes,
		}
		if h.UserID != nil {
			if u, ok := usersMap[*h.UserID]; ok {
				response[i].User = TaskUserResponse{
					ID:        u.ID,
					Name:      u.Name,
					Email:     u.Email,
					AvatarUrl: u.AvatarUrl,
				}
			} else {
				return nil, ErrDataIntegrity
			}
		}
	}

	return response, nil
}

func (s *service) HandleParticipantDelete(event events2.Event) error {
	data, ok := event.Data.(events2.ParticipantEvent)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	ctxT, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tasks, err := s.repo.GetByProjectIDAndUserID(ctxT, data.ProjectID, data.UserID)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		return nil
	}

	var histories []models.UpdateHistory
	now := time.Now()
	for i, t := range tasks {
		var changes []FieldChange
		oldJSON, err := json.Marshal(t)
		if err != nil {
			return err
		}

		if t.AssigneeID != nil && *t.AssigneeID == data.UserID {
			changes = append(changes, FieldChange{Field: "assignee_id", OldValue: t.AssigneeID})
			tasks[i].AssigneeID = nil
		}
		if t.ReporterID != nil && *t.ReporterID == data.UserID {
			changes = append(changes, FieldChange{Field: "reporter_id", OldValue: t.ReporterID})
			tasks[i].ReporterID = nil
		}

		if len(changes) == 0 {
			continue
		}

		changesJSON, err := json.Marshal(changes)
		if err != nil {
			return err
		}
		newJSON, err := json.Marshal(t)
		if err != nil {
			return err
		}

		histories = append(histories, models.UpdateHistory{
			CreatedAt: now,
			UserID:    nil,
			TaskID:    t.ID,
			Old:       oldJSON,
			New:       newJSON,
			Changes:   changesJSON,
		})
	}

	if err := s.repo.DeleteParticipantInTask(ctxT, tasks, histories); err != nil {
		log.Printf("err delete participant in tasks %s", err)
		return err
	}
	return nil
}

func equalUUIDPtr(a, b *uuid.UUID) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalUintPtr(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalStringPtr(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func equalTimePtr(a, b *time.Time) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Equal(*b)
}

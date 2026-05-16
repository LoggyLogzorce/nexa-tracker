package task

import (
	"errors"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
	"nexa-task-tracker/internal/pkg/nullable"
	"nexa-task-tracker/internal/pkg/response"
	"nexa-task-tracker/internal/pkg/validation"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=100"`
	Description *string    `json:"description"`
	StatusID    *uint      `json:"status_id"`
	PriorityID  *uint      `json:"priority_id"`
	AssigneeID  *uuid.UUID `json:"assignee_id" binding:"omitempty,uuid"`
	Deadline    *time.Time `json:"deadline" binding:"omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string         `json:"title" binding:"omitempty,min=1,max=100"`
	Description nullable.String `json:"description"`
	StatusID    nullable.Uint   `json:"status_id"`
	PriorityID  nullable.Uint   `json:"priority_id"`
	AssigneeID  nullable.UUID   `json:"assignee_id"`
	Deadline    nullable.Time   `json:"deadline"`
	IsArchive   nullable.Bool   `json:"is_archive"`
}

type Param struct {
	Archived  string
	UserParam string
}

func (h *Handler) Create(c *gin.Context) {
	pID := c.Param("id")
	projectID, err := uuid.Parse(pID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	rID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	reporterID := rID.(uuid.UUID)

	task := &Task{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		ProjectID:   projectID,
		StatusID:    req.StatusID,
		PriorityID:  req.PriorityID,
		AssigneeID:  req.AssigneeID,
		ReporterID:  &reporterID,
	}

	taskNew, err := h.service.Create(c.Request.Context(), task)
	if err != nil {
		switch {
		case errors.Is(err, ErrDataIntegrity):
			response.Error(c, http.StatusInternalServerError, "data integrity error")
		case errors.Is(err, ErrAssigneeNotInProject):
			response.Error(c, http.StatusBadRequest, "assignee not in project")
		case errors.Is(err, ErrPriorityNotInProject):
			response.Error(c, http.StatusBadRequest, "priority not in project")
		case errors.Is(err, ErrStatusNotInProject):
			response.Error(c, http.StatusBadRequest, "status not in project")
		case errors.Is(err, ErrInvalidAssigneeRole):
			response.Error(c, http.StatusBadRequest, "invalid assignee role")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to create task")
		}
		return
	}

	response.Success(c, http.StatusCreated, taskNew)
}

func (h *Handler) GetByID(c *gin.Context) {
	tId := c.Param("task_id")
	taskID, err := strconv.ParseUint(tId, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}
	archived := c.Query("archived")

	task, err := h.service.GetByID(c.Request.Context(), uint(taskID), archived)
	if err != nil {
		switch {
		case errors.Is(err, ErrTaskNotFound):
			response.Error(c, http.StatusNotFound, "task not found")
		case errors.Is(err, ErrStatusNotInProject):
			response.Error(c, http.StatusBadRequest, "status not in project")
		case errors.Is(err, ErrPriorityNotInProject):
			response.Error(c, http.StatusBadRequest, "priority not in project")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get task")
		}
		return
	}

	response.Success(c, http.StatusOK, task)
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	pID := c.Param("id")
	projectID, err := uuid.Parse(pID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}
	isArchive := c.Query("archived")
	param := Param{
		Archived: isArchive,
	}

	tasks, err := h.service.GetByProjectID(c.Request.Context(), projectID, param)
	if err != nil {
		switch {
		case errors.Is(err, ErrDataIntegrity):
			response.Error(c, http.StatusInternalServerError, "data integrity error")
			return
		case errors.Is(err, ErrAssigneeNotInProject):
			response.Error(c, http.StatusBadRequest, "assignee not in project")
			return
		case errors.Is(err, ErrPriorityNotInProject):
			response.Error(c, http.StatusBadRequest, "priority not in project")
			return
		case errors.Is(err, ErrStatusNotInProject):
			response.Error(c, http.StatusBadRequest, "status not in project")
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get tasks list")
		}
		return
	}
	response.Success(c, http.StatusOK, tasks)
}

func (h *Handler) GetByUserID(c *gin.Context) {
	uID, ok := c.Get(ctxkeys.UserIDKey)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	t := c.Query("type")
	param := Param{
		UserParam: t,
	}

	tasks, err := h.service.GetByUserID(c.Request.Context(), uID.(uuid.UUID), param)
	if err != nil {
		switch {
		case errors.Is(err, ErrDataIntegrity):
			response.Error(c, http.StatusInternalServerError, "data integrity error")
			return
		case errors.Is(err, ErrAssigneeNotInProject):
			response.Error(c, http.StatusBadRequest, "assignee not in project")
			return
		case errors.Is(err, ErrPriorityNotInProject):
			response.Error(c, http.StatusBadRequest, "priority not in project")
			return
		case errors.Is(err, ErrStatusNotInProject):
			response.Error(c, http.StatusBadRequest, "status not in project")
			return
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get tasks list")
		}
		return
	}
	response.Success(c, http.StatusOK, tasks)
}

func (h *Handler) Update(c *gin.Context) {
	tId := c.Param("task_id")
	taskID, err := strconv.ParseUint(tId, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	uID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := uID.(uuid.UUID)

	var req UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	archived := c.Query("archived")

	task, err := h.service.Update(c.Request.Context(), uint(taskID), &req, archived, userID)
	if err != nil {
		switch {
		case errors.Is(err, ErrTaskNotFound):
			response.Error(c, http.StatusNotFound, "task not found")
		case errors.Is(err, ErrStatusNotInProject):
			response.Error(c, http.StatusBadRequest, "status not in project")
		case errors.Is(err, ErrPriorityNotInProject):
			response.Error(c, http.StatusBadRequest, "priority not in project")
		case errors.Is(err, ErrAssigneeNotInProject):
			response.Error(c, http.StatusBadRequest, "assignee not in project")
		case errors.Is(err, ErrNoFieldsToUpdate):
			response.Error(c, http.StatusBadRequest, "no fields to update")
		case errors.Is(err, ErrInvalidAssigneeRole):
			response.Error(c, http.StatusBadRequest, "invalid assignee role")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update task")
		}
		return
	}
	response.Success(c, http.StatusOK, task)
}

func (h *Handler) Delete(c *gin.Context) {
	tId := c.Param("task_id")
	taskID, err := strconv.ParseUint(tId, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	uID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := uID.(uuid.UUID)

	if err := h.service.Delete(c.Request.Context(), uint(taskID), userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to delete task")
		return
	}

	response.Success(c, http.StatusOK, "task deleted")
}

func (h *Handler) GetHistoryByTaskID(c *gin.Context) {
	tId := c.Param("task_id")
	taskID, err := strconv.ParseUint(tId, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	history, err := h.service.GetHistoryByTaskID(c.Request.Context(), uint(taskID))
	if err != nil {
		switch {
		case errors.Is(err, ErrDataIntegrity):
			response.Error(c, http.StatusInternalServerError, "data integrity error")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get task history")
		}
		return
	}
	response.Success(c, http.StatusOK, history)
}

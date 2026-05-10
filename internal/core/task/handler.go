package task

import (
	"errors"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
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
	Title       *string `json:"title" binding:"omitempty,min=1,max=100"`
	Description *string `json:"description"`
	StatusID    *uint   `json:"status_id"`
	PriorityID  *uint   `json:"priority_id"`
	AssigneeID  *string `json:"assignee_id" binding:"omitempty,uuid"`
	Deadline    *string `json:"deadline"`
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
	task, err := h.service.GetByID(c.Request.Context(), uint(taskID))
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

	tasks, err := h.service.GetByProjectID(c.Request.Context(), projectID)
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
		default:
			response.Error(c, http.StatusInternalServerError, "failed to get tasks list")
		}
		return
	}
	response.Success(c, http.StatusOK, tasks)
}

func (h *Handler) Update(c *gin.Context) {
	// TODO: Update task
	c.JSON(http.StatusOK, gin.H{"message": "update task endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Delete task
	c.JSON(http.StatusOK, gin.H{"message": "delete task endpoint"})
}

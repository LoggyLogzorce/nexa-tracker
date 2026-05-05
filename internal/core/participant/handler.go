package participant

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/pkg/response"
	"nexa-task-tracker/internal/pkg/validation"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type AddParticipantRequest struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	Role   string `json:"role" binding:"required,oneof=owner member read_only"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=owner member read_only"`
}

func (h *Handler) AddParticipant(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid user id")
		return
	}

	participant := &ProjectParticipant{
		ProjectID: projectID,
		UserID:    userID,
		Role:      req.Role,
	}

	if err := h.service.AddParticipant(c.Request.Context(), participant); err != nil {
		switch {
		case errors.Is(err, ErrParticipantIDExists):
			response.Error(c, http.StatusConflict, "participant with this id already exists")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to add participant")
		}
		return
	}

	response.Success(c, http.StatusCreated, participant)
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	participants, err := h.service.GetByProjectID(c.Request.Context(), projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get participants")
		return
	}

	response.Success(c, http.StatusOK, participants)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	userID := c.Param("user_id")
	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	participant := &ProjectParticipant{
		ProjectID: projectID,
		UserID:    userIdUUID,
		Role:      req.Role,
	}

	if err := h.service.UpdateRole(c.Request.Context(), participant); err != nil {
		switch {
		case errors.Is(err, ErrParticipantNotFound):
			response.Error(c, http.StatusNotFound, "participant not found")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update role")
		}
		return
	}

	response.Success(c, http.StatusOK, participant)
}

func (h *Handler) RemoveParticipant(c *gin.Context) {
	idStr := c.Param("id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	userID := c.Param("user_id")
	userIdUUID, err := uuid.Parse(userID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	participant := &ProjectParticipant{
		ProjectID: projectID,
		UserID:    userIdUUID,
	}

	if err := h.service.RemoveParticipant(c.Request.Context(), participant); err != nil {
		switch {
		case errors.Is(err, ErrParticipantNotFound):
			response.Error(c, http.StatusNotFound, "participant not found")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to remove participant")
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "participant removed successfully"})
}

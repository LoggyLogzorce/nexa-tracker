package participant

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/pkg/response"
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
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	var req AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
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

	if err := h.service.AddParticipant(participant); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to add participant")
		return
	}

	response.Success(c, http.StatusCreated, participant)
}

func (h *Handler) GetByProjectID(c *gin.Context) {
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	participants, err := h.service.GetByProjectID(projectID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get participants")
		return
	}

	response.Success(c, http.StatusOK, participants)
}

func (h *Handler) UpdateRole(c *gin.Context) {
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	userID := c.Param("user_id")

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.UpdateRole(projectID, userID, req.Role); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to update role")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "role updated successfully"})
}

func (h *Handler) RemoveParticipant(c *gin.Context) {
	idStr := c.Param("project_id")
	projectID, err := uuid.Parse(idStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid project id")
		return
	}

	userID := c.Param("user_id")

	if err := h.service.RemoveParticipant(projectID, userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to remove participant")
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "participant removed successfully"})
}

package attachment

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"nexa-task-tracker/internal/ctxkeys"
	"nexa-task-tracker/internal/pkg/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	userIDVal, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}
	userID := userIDVal.(uuid.UUID)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "file is required")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to open file")
		return
	}
	defer file.Close()

	// защита от отсутствия Content-Length и лимит размера
	reader := io.LimitReader(file, 100<<20)

	attachment, err := h.service.Upload(
		c.Request.Context(),
		uint(taskID),
		userID,
		fileHeader.Filename,
		reader,
	)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to upload attachment")
		return
	}

	response.Success(c, http.StatusCreated, attachment)
}

func (h *Handler) GetByTaskID(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	attachments, err := h.service.GetByTaskID(c.Request.Context(), uint(taskID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get attachments")
		return
	}

	response.Success(c, http.StatusOK, attachments)
}

func (h *Handler) Download(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	aID := c.Param("attachment_id")
	attachmentID, err := strconv.ParseUint(aID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	attachment, err := h.service.GetByID(c.Request.Context(), uint(attachmentID), uint(taskID))
	if err != nil {
		if errors.Is(err, ErrAttachmentNotFound) {
			response.Error(c, http.StatusNotFound, "attachment not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "failed to get attachment")
		return
	}

	c.FileAttachment(attachment.FilePath, attachment.Filename)
}

func (h *Handler) Delete(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	aID := c.Param("attachment_id")
	attachmentID, err := strconv.ParseUint(aID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid attachment id")
		return
	}

	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	err = h.service.Delete(c.Request.Context(), uint(attachmentID), uint(taskID), userID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, ErrAttachmentNotFound):
			response.Error(c, http.StatusNotFound, "attachment not found")
		case errors.Is(err, ErrNotAttachmentOwner):
			response.Error(c, http.StatusForbidden, "user is not the attachment owner")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete attachment")
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "attachment deleted successfully"})
}

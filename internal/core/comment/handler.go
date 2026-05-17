package comment

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"nexa-task-tracker/internal/ctxkeys"
	"nexa-task-tracker/internal/pkg/response"
	"nexa-task-tracker/internal/pkg/validation"
	"strconv"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

type CommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000"`
}

func (h *Handler) Create(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	comment := &Comment{
		UserID:  userID.(uuid.UUID),
		TaskID:  uint(taskID),
		Content: req.Content,
	}

	resp, err := h.service.Create(c.Request.Context(), comment)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to create comment")
		return
	}
	response.Success(c, http.StatusCreated, resp)
}

func (h *Handler) GetByTaskID(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}
	comments, err := h.service.GetByTaskID(c.Request.Context(), uint(taskID))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "failed to get comments")
		return
	}

	response.Success(c, http.StatusOK, comments)
}

func (h *Handler) Update(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	cID := c.Param("comment_id")
	commentID, err := strconv.ParseUint(cID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		status, msg := validation.ParseError(err)
		response.Error(c, status, msg)
		return
	}

	comment := &Comment{
		ID:      uint(commentID),
		UserID:  userID.(uuid.UUID),
		TaskID:  uint(taskID),
		Content: req.Content,
	}

	resp, err := h.service.Update(c.Request.Context(), comment)
	if err != nil {
		switch {
		case errors.Is(err, ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "comment not found")
		case errors.Is(err, ErrNotCommentOwner):
			response.Error(c, http.StatusForbidden, "user is not the comment owner")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to update comment")
		}
		return
	}

	response.Success(c, http.StatusOK, resp)
}

func (h *Handler) Delete(c *gin.Context) {
	tID := c.Param("task_id")
	taskID, err := strconv.ParseUint(tID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	cID := c.Param("comment_id")
	commentID, err := strconv.ParseUint(cID, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid task id")
		return
	}

	userID, exists := c.Get(ctxkeys.UserIDKey)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "user not authenticated")
		return
	}

	if err := h.service.Delete(c.Request.Context(), uint(commentID), uint(taskID), userID.(uuid.UUID)); err != nil {
		switch {
		case errors.Is(err, ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, "comment not found")
		case errors.Is(err, ErrNotCommentOwner):
			response.Error(c, http.StatusForbidden, "user is not the comment owner")
		default:
			response.Error(c, http.StatusInternalServerError, "failed to delete comment")
		}
		return
	}

	response.Success(c, http.StatusOK, gin.H{"message": "comment deleted successfully"})
}

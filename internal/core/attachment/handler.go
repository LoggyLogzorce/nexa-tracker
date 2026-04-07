package attachment

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Upload(c *gin.Context) {
	// TODO: Implement file upload
	c.JSON(http.StatusOK, gin.H{"message": "upload attachment endpoint"})
}

func (h *Handler) GetByTaskID(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "get attachments by task endpoint"})
}

func (h *Handler) Download(c *gin.Context) {
	// TODO: Implement file download
	c.JSON(http.StatusOK, gin.H{"message": "download attachment endpoint"})
}

func (h *Handler) Delete(c *gin.Context) {
	// TODO: Implement
	c.JSON(http.StatusOK, gin.H{"message": "delete attachment endpoint"})
}

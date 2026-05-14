package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func BodySizeLimit(defaultLimit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := defaultLimit

		if c.FullPath() == "/api/v1/projects/:id/tasks/:task_id/attachments" &&
			c.Request.Method == http.MethodPost {
			limit = 100 << 20 // 100 MB
		}

		c.Request.Body = http.MaxBytesReader(
			c.Writer,
			c.Request.Body,
			limit,
		)

		c.Next()
	}
}

package notify

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"nexa-task-tracker/internal/db"
	"nexa-task-tracker/internal/middleware"
	"nexa-task-tracker/internal/pkg/events"
)

// TODO: Implement notification module
func Init(database *gorm.DB, eventBus *events.EventBus, engine *gin.Engine, secret string) {
	if err := db.Migrate(database, &Notification{}); err != nil {
		log.Fatal("Failed to run notify migrations:", err)
	}

	repo := NewRepository(database)
	svc := NewService(repo)
	hdl := NewHandler(svc)

	eventBus.Subscribe(events.UserDeleted, svc.HandleUserDeleted)

	// Notification routes
	notifications := engine.Group("/api/v1/notifications")
	notifications.Use(middleware.Auth(secret))
	{
		notifications.GET("", hdl.List)
		notifications.PUT("/:id/read", func(c *gin.Context) { c.JSON(200, gin.H{"message": "mark as read"}) })
		notifications.PUT("/read-all", func(c *gin.Context) { c.JSON(200, gin.H{"message": "mark all as read"}) })
	}
}

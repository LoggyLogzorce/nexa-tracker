package main

import (
	"log"
	_ "os"

	"github.com/joho/godotenv"
	"nexa-task-tracker/internal/api"
	"nexa-task-tracker/internal/config"
	"nexa-task-tracker/internal/core/attachment"
	"nexa-task-tracker/internal/core/auth"
	"nexa-task-tracker/internal/core/comment"
	"nexa-task-tracker/internal/core/history"
	"nexa-task-tracker/internal/core/participant"
	"nexa-task-tracker/internal/core/priority"
	"nexa-task-tracker/internal/core/project"
	"nexa-task-tracker/internal/core/status"
	"nexa-task-tracker/internal/core/task"
	"nexa-task-tracker/internal/core/user"
	"nexa-task-tracker/internal/db"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	database, err := db.Connect(db.DatabaseConfig{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := db.Migrate(database,
		&user.User{},
		&auth.RefreshToken{},
		&project.Project{},
		&participant.ProjectParticipant{},
		&status.Status{},
		&priority.Priority{},
		&task.Task{},
		&comment.Comment{},
		&history.UpdateHistory{},
		&attachment.Attachment{},
	); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	userRepo := user.NewRepository(database)
	authRepo := auth.NewRepository(database)

	// Initialize services
	authService := auth.NewService(authRepo, userRepo, cfg.JWT.Secret, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	// Initialize handlers
	authHandler := auth.NewHandler(authService, cfg.Cookie.Domain, cfg.Cookie.SameSite, cfg.JWT.AccessExpiry, cfg.JWT.RefreshExpiry)

	// Setup router
	h := api.Handlers{
		AuthHdl: authHandler,
		// ...
	}
	router := api.NewRouter(h)
	// notify.Init(database, router)
	engine := router.Setup()

	// Start server
	addr := cfg.Server.Host + ":" + cfg.Server.Port
	log.Printf("Starting server on %s", addr)

	if err := engine.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

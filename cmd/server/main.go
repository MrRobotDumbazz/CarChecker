package main

import (
	"car-status-backend/internal/config"
	"car-status-backend/internal/database"
	"car-status-backend/internal/server"
	"car-status-backend/internal/services"
	"log"
	"os"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := os.MkdirAll(cfg.Storage.UploadPath, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	db, err := database.Connect(*cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.RunMigrations(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	imageService := services.NewImageService(
		db,
		cfg.Storage.UploadPath,
		cfg.Storage.MaxFileSize,
		cfg.Storage.AllowedTypes,
	)

	predictionService := services.NewPredictionService(db)

	mlClient := services.NewMLClient(
		cfg.MLService.BaseURL,
		cfg.MLService.Timeout,
		cfg.MLService.APIKey,
	)

	var queueService *services.QueueService
	if cfg.Queue.Enabled {
		queueService = services.NewQueueService(
			cfg.Queue.Type,
			cfg.Queue.URL,
			cfg.Queue.QueueName,
			db,
		)
		log.Printf("Queue service enabled: %s", cfg.Queue.Type)
	} else {
		log.Println("Queue service disabled, using direct ML client calls")
	}

	handlers := server.NewHandlers(
		imageService,
		predictionService,
		mlClient,
		queueService,
		db,
	)

	srv := server.NewServer(cfg.Server.Host + ":" + cfg.Server.Port)
	srv.RegisterRoutes(handlers)

	go srv.GracefulShutdown()

	log.Printf("Car Status Backend starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Upload path: %s", cfg.Storage.UploadPath)
	log.Printf("ML Service URL: %s", cfg.MLService.BaseURL)
	log.Printf("Max file size: %d bytes", cfg.Storage.MaxFileSize)
	log.Printf("Allowed file types: %v", cfg.Storage.AllowedTypes)

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
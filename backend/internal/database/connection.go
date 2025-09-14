package database

import (
	"car-status-backend/internal/config"
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(cfg config.Config) (*DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(15 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db *DB) RunMigrations() error {
	migrationSQL := `
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица для хранения загруженных изображений
CREATE TABLE IF NOT EXISTS car_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename VARCHAR(255) NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Таблица для результатов предсказаний ML модели
CREATE TABLE IF NOT EXISTS predictions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_id UUID REFERENCES car_images(id) ON DELETE CASCADE,
    cleanliness_status VARCHAR(20) CHECK (cleanliness_status IN ('clean', 'dirty')),
    cleanliness_confidence DECIMAL(5,4) CHECK (cleanliness_confidence BETWEEN 0 AND 1),
    integrity_status VARCHAR(20) CHECK (integrity_status IN ('intact', 'damaged')),
    integrity_confidence DECIMAL(5,4) CHECK (integrity_confidence BETWEEN 0 AND 1),
    processing_time_ms INTEGER,
    ml_model_version VARCHAR(50),
    additional_data JSONB,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

-- Опциональная таблица для job queue (если не используем RabbitMQ/Kafka)
CREATE TABLE IF NOT EXISTS prediction_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    image_id UUID REFERENCES car_images(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    scheduled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание индексов только если они не существуют
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_car_images_uploaded_at') THEN
        CREATE INDEX idx_car_images_uploaded_at ON car_images(uploaded_at);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_predictions_image_id') THEN
        CREATE INDEX idx_predictions_image_id ON predictions(image_id);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_predictions_status') THEN
        CREATE INDEX idx_predictions_status ON predictions(status);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_predictions_created_at') THEN
        CREATE INDEX idx_predictions_created_at ON predictions(created_at);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_prediction_jobs_status') THEN
        CREATE INDEX idx_prediction_jobs_status ON prediction_jobs(status);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_prediction_jobs_scheduled_at') THEN
        CREATE INDEX idx_prediction_jobs_scheduled_at ON prediction_jobs(scheduled_at);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_car_images_filename') THEN
        CREATE INDEX idx_car_images_filename ON car_images(filename);
    END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_indexes WHERE indexname = 'idx_predictions_statuses') THEN
        CREATE INDEX idx_predictions_statuses ON predictions(cleanliness_status, integrity_status);
    END IF;
END
$$;
`

	_, err := db.Exec(migrationSQL)
	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.PingContext(ctx)
}
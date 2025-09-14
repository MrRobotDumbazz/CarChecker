-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Таблица для хранения загруженных изображений
CREATE TABLE car_images (
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
CREATE TABLE predictions (
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
CREATE TABLE prediction_jobs (
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

-- Индексы для оптимизации
CREATE INDEX idx_car_images_uploaded_at ON car_images(uploaded_at);
CREATE INDEX idx_predictions_image_id ON predictions(image_id);
CREATE INDEX idx_predictions_status ON predictions(status);
CREATE INDEX idx_predictions_created_at ON predictions(created_at);
CREATE INDEX idx_prediction_jobs_status ON prediction_jobs(status);
CREATE INDEX idx_prediction_jobs_scheduled_at ON prediction_jobs(scheduled_at);
CREATE INDEX idx_car_images_filename ON car_images(filename);
CREATE INDEX idx_predictions_statuses ON predictions(cleanliness_status, integrity_status);
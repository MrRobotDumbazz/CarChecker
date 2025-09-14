package services

import (
	"car-status-backend/internal/database"
	"car-status-backend/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type QueueService struct {
	queueType string
	conn      interface{}
	queueName string
	db        *database.DB
}

func NewQueueService(queueType, queueURL, queueName string, db *database.DB) *QueueService {
	return &QueueService{
		queueType: queueType,
		queueName: queueName,
		db:        db,
	}
}

func (q *QueueService) PublishPredictionJob(imageID uuid.UUID, imagePath string) error {
	switch q.queueType {
	case "rabbitmq":
		return q.publishToRabbitMQ(imageID, imagePath)
	case "kafka":
		return q.publishToKafka(imageID, imagePath)
	default:
		return q.publishToDB(imageID, imagePath)
	}
}

func (q *QueueService) publishToRabbitMQ(imageID uuid.UUID, imagePath string) error {
	job := models.PredictionJobRequest{
		ImageID:   imageID,
		ImagePath: imagePath,
		CreatedAt: time.Now(),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	return fmt.Errorf("RabbitMQ implementation not available yet, job data: %s", string(body))
}

func (q *QueueService) publishToKafka(imageID uuid.UUID, imagePath string) error {
	job := models.PredictionJobRequest{
		ImageID:   imageID,
		ImagePath: imagePath,
		CreatedAt: time.Now(),
	}

	body, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	return fmt.Errorf("Kafka implementation not available yet, job data: %s", string(body))
}

func (q *QueueService) publishToDB(imageID uuid.UUID, imagePath string) error {
	job := &models.PredictionJob{
		ID:          uuid.New(),
		ImageID:     imageID,
		Status:      string(models.JobStatusPending),
		RetryCount:  0,
		MaxRetries:  3,
		ScheduledAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	query := `
		INSERT INTO prediction_jobs (id, image_id, status, retry_count, max_retries, scheduled_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := q.db.Exec(query,
		job.ID,
		job.ImageID,
		job.Status,
		job.RetryCount,
		job.MaxRetries,
		job.ScheduledAt,
		job.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create prediction job: %w", err)
	}

	return nil
}

func (q *QueueService) GetPendingJobs(limit int) ([]models.PredictionJob, error) {
	var jobs []models.PredictionJob
	query := `
		SELECT id, image_id, status, retry_count, max_retries,
		       scheduled_at, started_at, completed_at, error_message, created_at
		FROM prediction_jobs
		WHERE status = 'pending' AND scheduled_at <= NOW()
		ORDER BY scheduled_at ASC
		LIMIT $1
	`

	err := q.db.Select(&jobs, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending jobs: %w", err)
	}

	return jobs, nil
}

func (q *QueueService) UpdateJobStatus(jobID uuid.UUID, status models.JobStatus, errorMessage string) error {
	now := time.Now()
	var query string
	var args []interface{}

	switch status {
	case models.JobStatusProcessing:
		query = `UPDATE prediction_jobs SET status = $1, started_at = $2 WHERE id = $3`
		args = []interface{}{status, now, jobID}
	case models.JobStatusCompleted:
		query = `UPDATE prediction_jobs SET status = $1, completed_at = $2 WHERE id = $3`
		args = []interface{}{status, now, jobID}
	case models.JobStatusFailed:
		query = `UPDATE prediction_jobs SET status = $1, error_message = $2, completed_at = $3 WHERE id = $4`
		args = []interface{}{status, errorMessage, now, jobID}
	default:
		query = `UPDATE prediction_jobs SET status = $1 WHERE id = $2`
		args = []interface{}{status, jobID}
	}

	_, err := q.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}

	return nil
}

func (q *QueueService) RetryJob(jobID uuid.UUID) error {
	query := `
		UPDATE prediction_jobs
		SET status = 'pending', retry_count = retry_count + 1, scheduled_at = NOW() + INTERVAL '5 minutes'
		WHERE id = $1 AND retry_count < max_retries
	`

	result, err := q.db.Exec(query, jobID)
	if err != nil {
		return fmt.Errorf("failed to retry job: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("job not found or max retries exceeded")
	}

	return nil
}

func (q *QueueService) GetJobStats() (map[string]interface{}, error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending,
			COUNT(CASE WHEN status = 'processing' THEN 1 END) as processing,
			COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
			COUNT(CASE WHEN status = 'failed' THEN 1 END) as failed
		FROM prediction_jobs
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`

	var stats struct {
		Total      int `db:"total"`
		Pending    int `db:"pending"`
		Processing int `db:"processing"`
		Completed  int `db:"completed"`
		Failed     int `db:"failed"`
	}

	err := q.db.Get(&stats, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get job stats: %w", err)
	}

	return map[string]interface{}{
		"total":      stats.Total,
		"pending":    stats.Pending,
		"processing": stats.Processing,
		"completed":  stats.Completed,
		"failed":     stats.Failed,
	}, nil
}
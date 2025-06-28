package services

import (
	"database/sql"
	"encoding/json"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/models"
)

type WebService struct {
	db *sql.DB
}

func NewWebService(db *sql.DB) *WebService {
	return &WebService{
		db: db,
	}
}

// it returns - total jobs, completed jobs, failed jobs and error
func (s *WebService) GetStats() (int, int, int, error) {
	query := `
	SELECT
		(SELECT COUNT(*) FROM jobs),
		(SELECT COUNT(*) FROM job_entries WHERE status = 'completed'),
		(SELECT COUNT(*) FROM job_entries WHERE status = 'permanently_failed');
	`

	var count, completedCount, failedCount int
	err := s.db.QueryRow(query).Scan(&count, &completedCount, &failedCount)
	return count, completedCount, failedCount, err
}

// to fetch all jobs
func (s *WebService) GetAllJobs() ([]*models.Job, error) {
	query := `
			SELECT id, cron_id, cron_expr, type, payload, created_at, updated_at
			FROM jobs
		`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	jobs := make([]*models.Job, 0)
	for rows.Next() {
		var job models.Job
		var jsonPayload string
		err := rows.Scan(&job.Id, &job.CronId, &job.CronExpr, &job.Type, &jsonPayload, &job.CreatedAt, &job.UpdatedAt)
		json.Unmarshal([]byte(jsonPayload), &job.Payload)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}
	return jobs, nil
}

// to fetch entries of a job
func (s *WebService) GetJobEntriesById(id int, status string) ([]*models.JobEntry, error) {
	query := `
			SELECT id, status, retries, output, error, scheduled_at, completed_at, updated_at
			FROM job_entries WHERE job_id = $1`

	args := []any{id}
	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	jobEntries := make([]*models.JobEntry, 0)
	for rows.Next() {
		var jobEntry models.JobEntry
		if err := rows.Scan(&jobEntry.Id, &jobEntry.Status, &jobEntry.Retries, &jobEntry.Output, &jobEntry.Error, &jobEntry.ScheduledAt, &jobEntry.CompletedAt, &jobEntry.UpdatedAt); err != nil {
			return nil, err
		}
		jobEntries = append(jobEntries, &jobEntry)
	}
	return jobEntries, nil
}

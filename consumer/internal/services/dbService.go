package services

import (
	"database/sql"
	"time"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
)

type DBService struct {
	db *sql.DB
}

func NewDBService(db *sql.DB) *DBService {
	return &DBService{
		db: db,
	}
}

// returns jobId, retries, error
func (s *DBService) getExistingJobEntry(job *models.Job, scheduledTime string) (*models.JobEntry, error) {
	query := `SELECT id, job_id, status, retries, output, error, scheduled_at, completed_at, updated_at
	FROM job_entries 
	WHERE job_id = $1 AND scheduled_at = $2 AND status != 'permanently_failed';
	`

	var jobEntry models.JobEntry
	err := s.db.QueryRow(query, job.Id, scheduledTime).Scan(&jobEntry.Id, &jobEntry.JobId, &jobEntry.Status, &jobEntry.Retries, &jobEntry.Output, &jobEntry.Error, &jobEntry.ScheduledAt, &jobEntry.CompletedAt, &jobEntry.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &jobEntry, nil
}

// create a new job entry
func (s *DBService) createNewJobEntry(job *models.Job, sTime string) (*models.JobEntry, error) {
	query := `INSERT INTO job_entries(job_id, status, scheduled_at)
			  VALUES ($1, $2, $3) RETURNING id;
	`

	scheduledTime, err := time.Parse("2006-01-02T15:04:05.000000-07:00", sTime)
	if err != nil {
		return nil, err
	}

	var jobEntryId int
	if err := s.db.QueryRow(query, job.Id, "running", sTime).Scan(&jobEntryId); err != nil {
		return nil, err
	}

	return &models.JobEntry{
		Id:          jobEntryId,
		JobId:       job.Id,
		Status:      "running",
		ScheduledAt: scheduledTime,
	}, nil
}

// register a job as completed
func (s *DBService) markJobAsCompleted(output string, jobEntry *models.JobEntry) {
	query := `UPDATE job_entries SET status = 'completed', output = $1, updated_at = $2 WHERE id = $3`
	s.db.Exec(query, output, time.Now(), jobEntry.Id)
}

// register a job as permanently failed
func (s *DBService) markJobAsPermanentlyFailed(jobEntry *models.JobEntry) {
	query := `UPDATE job_entries SET status = 'permanently_failed' WHERE id = $1`
	s.db.Exec(query, jobEntry.Id)
}

// register a job as failed
func (s *DBService) markJobAsFailed(err error, updatedRetries int, jobEntry *models.JobEntry) {
	query := `UPDATE job_entries SET status = 'failed', error = $1, retries = $2, updated_at = $3 WHERE id = $4`
	s.db.Exec(query, err.Error(), updatedRetries, time.Now(), jobEntry.Id)
}

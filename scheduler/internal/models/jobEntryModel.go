package models

import (
	"database/sql"
	"time"
)

type JobEntry struct {
	Id          int            `json:"id"`
	JobId       int            `json:"job_id"`
	Status      string         `json:"status"`
	Retries     int            `json:"retries"`
	Output      sql.NullString `json:"output"`
	Error       sql.NullString `json:"error"`
	ScheduledAt time.Time      `json:"scheduled_at"`
	CompletedAt time.Time      `json:"completed_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

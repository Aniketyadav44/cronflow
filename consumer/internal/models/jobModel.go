package models

import "time"

type Job struct {
	Id        int            `json:"id"`
	CronId    int            `json:"cron_id"`
	CronExpr  string         `json:"cron_expr"`
	Type      string         `json:"type"`
	Payload   map[string]any `json:"payload"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

package config

import (
	"database/sql"

	"github.com/robfig/cron/v3"
)

func loadCron(db *sql.DB) (*cron.Cron, error) {
	cron := cron.New()
	cron.Start()

	return cron, nil

	// // starting cron jobs from db
	// query := `SELECT id, cron_id, cron_expr, payload from jobs`
	// rows, err := db.Query(query)
	// if err != nil {
	// 	return nil, err
	// }

	// jobs := make([]*models.Job, 0)
	// for rows.Next() {
	// 	var job *models.Job
	// 	if err := rows.Scan(&job.Id, &job.CronId, &job.CronExpr, &job.Payload); err != nil {
	// 		return nil, err
	// 	}
	// 	jobs = append(jobs, job)
	// }

	// for _, v := range jobs{

	// }

}

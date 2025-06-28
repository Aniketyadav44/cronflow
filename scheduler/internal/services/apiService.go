package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/models"
	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
)

type ApiService struct {
	db        *sql.DB
	cron      *cron.Cron
	rabbitmq  *amqp091.Connection
	mqChannel *amqp091.Channel
}

func NewApiService(db *sql.DB, cron *cron.Cron, rmq *amqp091.Connection, mqChannel *amqp091.Channel) *ApiService {
	return &ApiService{
		db:        db,
		cron:      cron,
		rabbitmq:  rmq,
		mqChannel: mqChannel,
	}
}

// only for testing. to get all list of actually scheduled cron jobs on machine
func (s *ApiService) GetMachineCrons() []map[string]any {
	cronEntries := s.cron.Entries()
	entries := make([]map[string]any, 0)

	for _, v := range cronEntries {
		e := map[string]any{
			"cron_id": int(v.ID),
			"next":    v.Next,
		}
		entries = append(entries, e)
	}
	return entries
}

func (s *ApiService) CreateNewJob(job *models.Job) error {
	// first putting this job in db
	query := `INSERT INTO jobs(cron_expr, type, payload)
			  VALUES ($1, $2, $3) RETURNING id;
	`
	payloadJSON, _ := json.Marshal(job.Payload)
	var jobId int
	if err := s.db.QueryRow(query, job.CronExpr, job.Type, payloadJSON).Scan(&jobId); err != nil {
		return err
	}
	job.Id = jobId

	// scheduling a cron job
	id, err := s.cron.AddFunc(job.CronExpr, func() {
		log.Println("running cron job: publishing to rabbitmq")
		q, err := s.mqChannel.QueueDeclare("cron_events", false, false, false, false, nil)
		if err != nil {
			log.Println("failed creating a queue for rabbitmq: ", err.Error())
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		jsonBody, err := json.Marshal(map[string]any{
			"job":  job,
			"time": time.Now().Format("2006-01-02T15:04:05.000000-07:00"),
		})
		if err != nil {
			log.Println("failed creating payload for rabbitmq: ", err.Error())
			return
		}

		if err := s.mqChannel.PublishWithContext(ctx, "", q.Name, false, false, amqp091.Publishing{
			ContentType: "application/json",
			Body:        jsonBody,
		}); err != nil {
			log.Println("failed publishing to rabbitmq: ", err.Error())
			return
		}
		log.Println("event published to rabbitmq!")

	})
	if err != nil {
		// on cron scheduling error, deleting the job created in db
		delQuery := `DELETE FROM jobs WHERE id = $1`
		s.db.Exec(delQuery, jobId)
		return err
	}

	updateQuery := `UPDATE jobs SET cron_id = $1 WHERE id = $2`
	s.db.Exec(updateQuery, id, jobId)
	job.CronId = int(id)
	return nil
}

// To delete a job and all of it's entries
func (s *ApiService) DeleteJob(id int) error {
	query := `SELECT cron_id from jobs WHERE id = $1`

	var cronId int
	s.db.QueryRow(query, id).Scan(&cronId)

	s.cron.Remove(cron.EntryID(cronId))

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`DELETE FROM job_entries WHERE job_id = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM jobs WHERE id = $1`, id); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
	"github.com/rabbitmq/amqp091-go"
)

const MaxJobRetries = 3

type RMQService struct {
	dbService *DBService
	conn      *amqp091.Connection
	channel   *amqp091.Channel
}

func NewRMQService(db *sql.DB, conn *amqp091.Connection, channel *amqp091.Channel) *RMQService {
	return &RMQService{
		dbService: NewDBService(db),
		conn:      conn,
		channel:   channel,
	}
}

func (s *RMQService) Start(ctx context.Context) {
	q, err := s.channel.QueueDeclare("cron_events", false, false, false, false, nil)
	if err != nil {
		log.Println("error in creating rabbitmq queue: ", err.Error())
		return
	}

	msgs, err := s.channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Println("error in creating a consume channel for rabbitmq: ", err.Error())
		return
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Stopping rabbitmq...")
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Println("RabbitMQ message channel is closed.")
					return
				}
				processMessage(&msg, s.dbService)
			}
		}
	}()
	log.Println("RabbitMQ Consumer Running: Waiting for messages...")

	<-ctx.Done()
}

// function to process the queue's message.
// on errors, will NACK() - to let this msg get received again
// on success, will ACK() the message
func processMessage(msg *amqp091.Delivery, dbService *DBService) {
	log.Println("Received message on RabbitMQ channel: ", string(msg.Body))

	// parsing message body which has keys "job"[Job json] and "time"[Schedule time string]
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		log.Println("error in extracting message payload: ", err.Error())
		msg.Ack(false)
		return
	}

	// parsing job json from the message body json
	jobBody, ok := body["job"].(map[string]any)
	if !ok {
		log.Println("error in extracting job json: ", body["job"])
		msg.Ack(false)
		return
	}
	// parsing scheduled time from the message body json
	sTime, ok := body["time"].(string)
	if !ok {
		log.Println("error in extracting scheduled time: ", body["time"])
		msg.Ack(false)
		return
	}

	// converting the job json to bytes, to convert it to models.Job
	jobBodyByte, err := json.Marshal(jobBody)
	if err != nil {
		log.Println("error in parsing job json: ", err.Error())
		msg.Ack(false)
		return
	}
	var job *models.Job
	if err := json.Unmarshal(jobBodyByte, &job); err != nil {
		log.Println("error in extracting job: ", err.Error())
		msg.Ack(false)
		return
	}

	// get retries of any existing job entry for this job id, scheduled time which was not failed
	var jobEntry *models.JobEntry
	j, err := dbService.getExistingJobEntry(job, sTime)
	if err != nil {
		log.Println("error in getting a job entry: ", err.Error())
		msg.Ack(false)
		return
	}
	if j != nil {
		// if a job entry already exists, update that in jobEntry variable
		jobEntry = j
	} else {
		createdEntry, err := dbService.createNewJobEntry(job, sTime)
		if err != nil {
			log.Println("error creating new entry in db: ", err.Error())
			msg.Ack(false)
			return
		}
		jobEntry = createdEntry
	}

	// checking if the retry count reached max retries
	if jobEntry.Retries >= MaxJobRetries {
		dbService.markJobAsPermanentlyFailed(jobEntry)
		msg.Ack(false)
		return
	}

	switch job.Type {
	case "ping":
		if err := processPingJob(dbService, job, jobEntry, sTime); err != nil {
			handleJobError(dbService, err, msg, jobEntry)
			return
		}
	case "email":
		if err := processEmailJob(dbService, job, jobEntry, sTime); err != nil {
			handleJobError(dbService, err, msg, jobEntry)
			return
		}
	case "slack":
		if err := processSlackJob(dbService, job, jobEntry, sTime); err != nil {
			handleJobError(dbService, err, msg, jobEntry)
			return
		}
	case "webhook":
		if err := processWebhookJob(dbService, job, jobEntry, sTime); err != nil {
			handleJobError(dbService, err, msg, jobEntry)
			return
		}
	default:
		handleJobError(dbService, fmt.Errorf("invalid event type: %s", job.Type), msg, jobEntry)
	}

	msg.Ack(false)
}

// a function for handling errors. It sleeps for 2 second and unacknowledges the message
// and returns the message to queue
// It will retry for [MaxJobRetries] times and then will ack the message and log with error in db.
func handleJobError(dbService *DBService, err error, msg *amqp091.Delivery, jobEntry *models.JobEntry) {
	log.Println("error in processing job: ", err.Error(), ", retries: ", jobEntry.Retries)
	time.Sleep(2 * time.Second)
	dbService.markJobAsFailed(err, jobEntry.Retries+1, jobEntry)
	msg.Nack(false, true)
}

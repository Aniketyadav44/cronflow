package services

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
)

func processSlackJob(dbService *DBService, job *models.Job, jobEntry *models.JobEntry, sTime string) error {
	log.Println("Received a slack event")

	url, ok := job.Payload["url"].(string)
	if !ok {
		return fmt.Errorf("invalid url: %s", job.Payload["url"])
	}

	msg, ok := job.Payload["mg"].(string)
	if !ok {
		return fmt.Errorf("invalid msg: %s", job.Payload["msg"])
	}

	body := fmt.Sprintf("{\"text\": \"%s\"}", msg)

	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("error sending message: %s", err.Error())
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("non 200/201 status code: %d", res.StatusCode)
	} else {
		dbService.markJobAsCompleted("Message sent successfully", jobEntry)
	}

	return nil
}

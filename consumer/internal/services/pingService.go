package services

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
)

func processPingJob(dbService *DBService, job *models.Job, jobEntry *models.JobEntry, sTime string) error {
	log.Println("Received a ping event")

	url, ok := job.Payload["url"].(string)
	if !ok {
		return fmt.Errorf("invalid url: %s", job.Payload["url"])
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		log.Println("Ping event executed successfully")
		// success api call, enter in db
		dbService.markJobAsCompleted("URL pinged successfully!", jobEntry)
	} else {
		// failed api call
		return fmt.Errorf("url ping resulted %d code", resp.StatusCode)
	}

	return nil
}

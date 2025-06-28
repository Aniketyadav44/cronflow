package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
)

func processWebhookJob(dbService *DBService, job *models.Job, jobEntry *models.JobEntry, sTime string) error {
	log.Println("Received a webhook job")

	url, ok := job.Payload["url"].(string)
	if !ok {
		return fmt.Errorf("invalid webhook url: %s", job.Payload["url"])
	}

	body, ok := job.Payload["body"].(string)
	if !ok {
		return fmt.Errorf("invalid body: %s", job.Payload["body"])
	}

	res, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return fmt.Errorf("error calling POST: %s", err.Error())
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("non 200/201 status code in post api: %d", res.StatusCode)
	} else {
		defer res.Body.Close()
		resBody, _ := io.ReadAll(res.Body)
		dbService.markJobAsCompleted("API called successfully: "+string(resBody), jobEntry)
	}
	return nil
}

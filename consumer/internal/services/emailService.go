package services

import (
	"fmt"
	"log"
	"os"

	"github.com/Aniketyadav44/cronflow/consumer/internal/models"
	"github.com/go-gomail/gomail"
)

func processEmailJob(dbServicedb *DBService, job *models.Job, jobEntry *models.JobEntry, sTime string) error {
	log.Println("Received an email event")

	email, ok := job.Payload["email"].(string)
	if !ok {
		return fmt.Errorf("invalid email: %s", job.Payload["email"])
	}
	subject, ok := job.Payload["subject"].(string)
	if !ok {
		return fmt.Errorf("invalid email subject: %s", job.Payload)
	}
	body, ok := job.Payload["body"].(string)
	if !ok {
		return fmt.Errorf("invalid email body: %s", job.Payload["body"])
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "bot@demomailtrap.co")
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	emailPass := os.Getenv("MAIL_PASS")
	dialer := gomail.NewDialer("live.smtp.mailtrap.io", 587, "api", emailPass)

	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("error sending email: %s", err.Error())
	} else {
		dbServicedb.markJobAsCompleted("Email sent successfully", jobEntry)
	}
	return nil
}

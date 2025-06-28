package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/models"
	"github.com/Aniketyadav44/cronflow/scheduler/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type ApiHandler struct {
	service *services.ApiService
}

func NewApiHandler(service *services.ApiService) *ApiHandler {
	return &ApiHandler{
		service: service,
	}
}

// only for testing. to get all list of actually scheduled cron jobs on machine
func (h *ApiHandler) GetMachineCrons(c *gin.Context) {
	entries := h.service.GetMachineCrons()
	c.JSON(http.StatusOK, gin.H{"crons": entries})
}

// This schedules a cron job and makes entry to postgres db
func (h *ApiHandler) CreateNewJob(c *gin.Context) {
	hour := c.PostForm("hour")
	minute := c.PostForm("minute")
	payload := make(map[string]any)
	payload["hour"] = hour
	payload["minute"] = minute
	taskType := c.PostForm("type")
	switch taskType {
	case "ping":
		payload["url"] = c.PostForm("url")
	case "email":
		payload["email"] = c.PostForm("email")
		payload["subject"] = c.PostForm("subject")
		payload["body"] = c.PostForm("body")
	case "slack":
		payload["url"] = c.PostForm("url")
		payload["msg"] = c.PostForm("msg")
	case "webhook":
		payload["url"] = c.PostForm("url")
		payload["body"] = c.PostForm("body")
	}

	job := &models.Job{
		CronExpr: fmt.Sprintf("%s %s * * *", minute, hour),
		Type:     taskType,
		Payload:  payload,
	}

	if err := h.service.CreateNewJob(job); err != nil {
		log.Println("error in creating cron job: ", err.Error())
		session := sessions.Default(c) // using sessions only for handling & showing errors
		session.Set("error", err.Error())
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/jobs/new")
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/")
}

// This deletes a scheduled job and from db as well
func (h *ApiHandler) DeleteJob(c *gin.Context) {
	id := c.Request.URL.Query().Get("id")
	if id == "" {
		c.String(http.StatusBadRequest, "invalid id query")
		return
	}

	idInt, _ := strconv.Atoi(id)
	err := h.service.DeleteJob(idInt)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		return
	}

	c.Redirect(http.StatusMovedPermanently, "/jobs")
}

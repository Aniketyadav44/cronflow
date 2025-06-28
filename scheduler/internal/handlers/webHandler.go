package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type WebHandler struct {
	service *services.WebService
}

func NewWebHandler(service *services.WebService) *WebHandler {
	return &WebHandler{
		service: service,
	}
}

// for root dashboard route "/"
func (h *WebHandler) Home(c *gin.Context) {
	total, completedCount, failedCount, err := h.service.GetStats()
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error occurred in fetching stats: %s", err.Error()))
		return
	}
	c.HTML(http.StatusOK, "home.html", gin.H{
		"Total":          total,
		"CompletedCount": completedCount,
		"FailedCount":    failedCount,
	})
}

// for create form page route "/jobs/new"
func (h *WebHandler) CreateNewJob(c *gin.Context) {
	session := sessions.Default(c)
	err := session.Get("error")
	if err != nil {
		session.Delete("error")
		session.Save()
	}
	c.HTML(http.StatusOK, "create-form.html", gin.H{
		"Hours":        24,
		"Minutes":      60,
		"ErrorMessage": err,
	})
}

// for listing all jobs route "/jobs"
func (h *WebHandler) ListJobs(c *gin.Context) {
	jobs, err := h.service.GetAllJobs()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "list-jobs.html", gin.H{
		"Jobs": jobs,
	})
}

// for a job entry route "/job-entries?id="
func (h *WebHandler) ListJobEntries(c *gin.Context) {
	id := c.Request.URL.Query().Get("id")
	status := c.Request.URL.Query().Get("status")
	if id == "" {
		c.String(http.StatusBadRequest, "Missing id query")
		return
	}
	idInt, _ := strconv.Atoi(id)
	jobEntries, err := h.service.GetJobEntriesById(idInt, status)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.HTML(http.StatusOK, "job-entries.html", gin.H{"JobId": id, "Entries": jobEntries, "Filter": status})
}

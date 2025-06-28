package web

import (
	"database/sql"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/handlers"
	"github.com/Aniketyadav44/cronflow/scheduler/internal/services"
	"github.com/gin-gonic/gin"
)

func registerWebRoutes(router *gin.Engine, db *sql.DB) {
	webService := services.NewWebService(db)
	webHandler := handlers.NewWebHandler(webService)

	d := router.Group("/")
	{
		d.GET("", webHandler.Home)
		d.GET("/jobs/new", webHandler.CreateNewJob)
		d.GET("/jobs", webHandler.ListJobs)
		d.GET("/job-entries", webHandler.ListJobEntries)
	}
}

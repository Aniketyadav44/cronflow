package web

import (
	"database/sql"

	"github.com/Aniketyadav44/cronflow/scheduler/internal/handlers"
	"github.com/Aniketyadav44/cronflow/scheduler/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
)

func registerApiRoutes(router *gin.Engine, db *sql.DB, cron *cron.Cron, rmq *amqp091.Connection, mqChannel *amqp091.Channel) {
	apiService := services.NewApiService(db, cron, rmq, mqChannel)
	apiHandler := handlers.NewApiHandler(apiService)

	api := router.Group("/api")
	{
		api.POST("/create", apiHandler.CreateNewJob)
		api.GET("/machine-crons", apiHandler.GetMachineCrons)
		api.POST("/delete-job", apiHandler.DeleteJob)
	}
}

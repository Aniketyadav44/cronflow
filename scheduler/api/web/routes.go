package web

import (
	"github.com/Aniketyadav44/cronflow/scheduler/internal/config"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, cfg *config.Config) {
	registerWebRoutes(router, cfg.Db)
	registerApiRoutes(router, cfg.Db, cfg.Cron, cfg.RabbitMQ, cfg.MQChannel)
}

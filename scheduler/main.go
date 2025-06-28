package main

import (
	"log"

	web "github.com/Aniketyadav44/cronflow/scheduler/api/web"
	"github.com/Aniketyadav44/cronflow/scheduler/internal/config"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error loading config: ", err.Error())
	}

	defer cfg.Db.Close()
	defer cfg.RabbitMQ.Close()
	defer cfg.Cron.Stop()
	defer cfg.MQChannel.Close()

	router := gin.Default()
	router.LoadHTMLGlob("api/web/templates/*.html")
	router.Static("/static", "./api/web/static")

	// for handling web interface's session
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	web.RegisterRoutes(router, cfg)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("error in starting server: ", err.Error())
	}
}

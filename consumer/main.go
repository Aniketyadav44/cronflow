package main

import (
	"context"
	"log"

	"github.com/Aniketyadav44/cronflow/consumer/internal/config"
	"github.com/Aniketyadav44/cronflow/consumer/internal/services"
)

const MAX_JOB_RETRIES = 3

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("error in loading config: ", err.Error())
	}

	defer cfg.Db.Close()
	defer cfg.RabbitMQ.Close()
	defer cfg.MQChannel.Close()

	rabbitmqService := services.NewRMQPService(cfg.Db, cfg.RabbitMQ, cfg.MQChannel)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rabbitmqService.Start(ctx)
}

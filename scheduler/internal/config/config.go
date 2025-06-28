package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/robfig/cron/v3"
)

type Config struct {
	Port      string
	Db        *sql.DB
	RabbitMQ  *amqp091.Connection
	MQChannel *amqp091.Channel
	Cron      *cron.Cron
}

func LoadConfig() (*Config, error) {
	// if err := godotenv.Load("../.env"); err != nil {
	// 	return nil, err
	// }

	webPort := os.Getenv("DASHBOARD_PORT")
	pgHost := os.Getenv("PG_HOST")
	pgPort := os.Getenv("PG_PORT")
	pgUser := os.Getenv("PG_USER")
	pgPass := os.Getenv("PG_PASSWORD")
	pgDbName := os.Getenv("PG_DBNAME")
	rmqUser := os.Getenv("RMQ_USER")
	rmqPass := os.Getenv("RMQ_PASSWORD")
	rmqHost := os.Getenv("RMQ_HOST")
	rmqPort := os.Getenv("RMQ_PORT")
	rmqVHost := os.Getenv("RMQ_VHOST")

	// loading postgres database
	db, err := loadDb(pgHost, pgPort, pgUser, pgPass, pgDbName)
	if err != nil {
		return nil, fmt.Errorf("error in loading db: %s", err.Error())
	}

	// loading rabbitmq
	rabbitmq, mqCh, err := loadRabbitMQ(rmqUser, rmqPass, rmqHost, rmqPort, rmqVHost)
	if err != nil {
		return nil, fmt.Errorf("error in loading rabbitMQ: %s", err.Error())
	}

	// loading cron
	cron, err := loadCron(db)
	if err != nil {
		return nil, fmt.Errorf("error in loading cron: %s", err.Error())
	}

	return &Config{
		Port:      webPort,
		Db:        db,
		RabbitMQ:  rabbitmq,
		MQChannel: mqCh,
		Cron:      cron,
	}, nil

}

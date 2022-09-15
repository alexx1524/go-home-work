package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	cfg "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/messaging/rabbitmq"
	internalscheduler "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := cfg.NewConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logger, err := internallogger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		log.Fatalln(err)
	}

	connectionString := config.Rabbit.ConnectionString
	exchangeName := config.Rabbit.Exchange
	queueName := config.Rabbit.Queue

	producer, err := rabbitmq.NewConnector(connectionString, exchangeName, queueName)
	if err != nil {
		log.Fatalln(err)
	}
	defer func(connector rabbitmq.Connector) {
		err := connector.Close()
		if err != nil {
			logger.Debug(fmt.Sprintf("Producer stopping error: %s", err.Error()))
		}
	}(producer)

	storage, err := sqlstorage.New(config.DBStorage.ConnectionString)
	if err != nil {
		log.Panicln(err)
	}
	defer storage.Close()
	logger.Info("Storage is created")

	scheduler := internalscheduler.NewScheduler(logger, storage)

	go func() {
		for event := range scheduler.GetNotificationChannel() {
			bytes, err := json.Marshal(event)
			if err != nil {
				logger.Error(fmt.Sprintf("Marshaling error %s", err.Error()))
			}

			err = producer.Send(bytes)
			if err != nil {
				logger.Error(fmt.Sprintf("Sending notification error %s", err.Error()))
			}
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	err = scheduler.Run(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("Scheduler starting error %s", err.Error()))
	}

	<-ctx.Done()
}

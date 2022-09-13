package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	cfg "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/messaging/rabbitmq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/sender_config.toml", "Path to configuration file")
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	connectionString := config.Rabbit.ConnectionString
	exchangeName := config.Rabbit.Exchange
	queueName := config.Rabbit.Queue

	consumer, err := rabbitmq.NewConnector(connectionString, exchangeName, queueName)
	if err != nil {
		logger.Error(err.Error())
	}
	defer func(consumer rabbitmq.Connector) {
		err := consumer.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}(consumer)

	channel, err := consumer.Consume()
	if err != nil {
		logger.Error(err.Error())
	}

	go func() {
		for event := range channel {
			message := fmt.Sprintf("Received message: %s", string(event.Body))
			logger.Debug(message)

			err = event.Ack(false)
			if err != nil {
				logger.Error(err.Error())
			}
		}
	}()

	<-ctx.Done()
}

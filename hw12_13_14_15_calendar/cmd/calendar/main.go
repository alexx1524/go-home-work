package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/http"
	internalstorage "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage/sql"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := cfg.NewConfig(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	logg, err := internallogger.New(config.Logger.Level, config.Logger.File)
	if err != nil {
		log.Fatalln(err)
	}

	storage, err := initStorage(config)
	if err != nil {
		log.Fatalln(err)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initStorage(config cfg.Config) (app.Storage, error) {
	switch config.StorageSource {
	case "in-memory":
		return memorystorage.New(), nil
	case "sql":
		return sqlstorage.New(config.DBStorage.ConnectionString)
	default:
		return nil, internalstorage.ErrorSourceUnsupported
	}
}

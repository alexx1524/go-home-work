package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/app"
	cfg "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/grpc"
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

	httpServer := internalhttp.NewServer(logg, calendar, config)
	grpcServer := internalgrpc.NewServer(logg, calendar, config)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running...")

	go func() {
		logg.Info("Http server is starting...")
		if err := httpServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error(fmt.Sprintf("listen: %s", err.Error()))
		}
	}()

	go func() {
		logg.Info("gRPC server is starting...")
		if err := grpcServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logg.Error(fmt.Sprintf("listen: %s", err.Error()))
		}
	}()

	<-ctx.Done()

	if err := httpServer.Stop(ctx); err != nil {
		logg.Error(fmt.Sprintf("Http server stopping error: %s", err.Error()))
	}

	if err := grpcServer.Stop(ctx); err != nil {
		logg.Error(fmt.Sprintf("gRPC server stopping error: %s", err.Error()))
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

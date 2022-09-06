package grpc

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	internallogger "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/grpc/eventpb"
	memorystorage "github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGrpcServer(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		log.Println(err)
	}
	defer os.RemoveAll(dir)
	file := filepath.Join(dir, "log.txt")

	cfg := config.Config{
		GRPCServer: config.GRPCServer{
			Host: "localhost",
			Port: "50021",
		},
	}

	logger, err := internallogger.New("debug", file)
	if err != nil {
		log.Println(err)
	}

	calendar := app.New(logger, memorystorage.New())

	t.Run("If event doesn't exist GetEventById returns error", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		server := NewServer(logger, calendar, cfg)

		go func() {
			if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}()

		// create gRpc client
		conn, err := grpc.DialContext(ctx, "localhost:50021",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := eventpb.NewCalendarClient(conn)
		e, err := client.GetEventByID(ctx, &eventpb.GetEventRequest{
			ID: "f20d0b91-7c04-481d-ab64-b1fca8f2e47f",
		})

		require.Nil(t, e)
		require.Error(t, err)

		err = server.Stop(ctx)
		cancel()

		if err != nil {
			log.Fatalln(err)
		}
	})

	t.Run("Create and get a new event", func(t *testing.T) {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)

		server := NewServer(logger, calendar, cfg)

		go func() {
			if err := server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln(err)
			}
		}()

		conn, err := grpc.DialContext(ctx, "localhost:50021",
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := eventpb.NewCalendarClient(conn)

		id := uuid.New().String()
		e, err := client.CreateEvent(ctx, &eventpb.CreateEventRequest{
			Event: &eventpb.Event{
				ID:          id,
				Title:       "title",
				Description: "description",
				UserID:      uuid.New().String(),
				StartDate:   timestamppb.New(time.Now()),
				EndDate:     timestamppb.New(time.Now().AddDate(0, 0, 1)),
			},
		})

		require.NotNil(t, e)
		require.NoError(t, err)

		e1, err := client.GetEventByID(ctx, &eventpb.GetEventRequest{
			ID: id,
		})

		require.NotNil(t, e1)
		require.NoError(t, err)

		err = server.Stop(ctx)
		cancel()

		if err != nil {
			log.Fatalln(err)
		}
	})
}

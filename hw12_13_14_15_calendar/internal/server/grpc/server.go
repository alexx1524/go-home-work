package grpc

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/codes"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/server/grpc/eventpb"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Server struct {
	app    Application
	config config.Config
	logger Logger
	server *grpc.Server
}

type Application interface {
	InsertEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEventByID(ctx context.Context, eventID uuid.UUID) (storage.Event, error)
	GetEventsCount(ctx context.Context) (int, error)
	GetEventsForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]storage.Event, error)
}

type Logger interface {
	Error(msg string)
	Info(msg string)
	LogGRPCRequest(code codes.Code, method, address string, requestDuration time.Duration)
}

type CalendarService struct {
	app    Application
	logger Logger
	eventpb.UnimplementedCalendarServer
}

func NewServer(logger Logger, app Application, config config.Config) *Server {
	interceptor := grpc.ChainUnaryInterceptor(
		LoggingInterceptor(logger),
	)
	grpcServer := grpc.NewServer(interceptor)
	service := &CalendarService{
		app:    app,
		logger: logger,
	}

	eventpb.RegisterCalendarServer(grpcServer, service)

	server := &Server{
		app:    app,
		config: config,
		logger: logger,
		server: grpcServer,
	}

	return server
}

func (s *Server) Start(ctx context.Context) error {
	endpoint := fmt.Sprintf("%v:%v", s.config.GRPCServer.Host, s.config.GRPCServer.Port)
	lsn, err := net.Listen("tcp", endpoint)
	if err != nil {
		s.logger.Error(fmt.Sprintf("fail start gprc server: %s", err.Error()))
	}

	if err := s.server.Serve(lsn); err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.logger.Error(fmt.Sprintf("listen: %s", err.Error()))
		os.Exit(1)
	}

	<-ctx.Done()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("grpc is stopping...")
	s.server.GracefulStop()

	<-ctx.Done()

	return nil
}

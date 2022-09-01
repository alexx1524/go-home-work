package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EventsSearchCriteria struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type Server struct {
	logger Logger
	app    Application
	server *http.Server
	Router *mux.Router
}

type Logger interface {
	Error(msg string)
	Info(msg string)
	LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration)
}

type Application interface {
	InsertEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEventByID(ctx context.Context, eventID uuid.UUID) (storage.Event, error)
	GetEventsCount(ctx context.Context) (int, error)
	GetEventsForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]storage.Event, error)
}

type loggingMiddleware struct {
	logger Logger
}

func NewServer(logger Logger, app Application, config config.Config) *Server {
	router := mux.NewRouter()
	httpServer := &http.Server{
		Addr:              net.JoinHostPort(config.HTTPServer.Host, config.HTTPServer.Port),
		ReadHeaderTimeout: time.Duration(config.HTTPServer.ReadHeaderTimeoutSeconds) * time.Second,
		Handler:           router,
	}

	server := &Server{
		logger: logger,
		app:    app,
		server: httpServer,
		Router: router,
	}

	server.InitializeEventsRoutes()

	loggingMiddleware := loggingMiddleware{
		logger: logger,
	}

	router.Use(loggingMiddleware.Process)

	return server
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Error(err.Error())
		log.Fatalln("Ошибка запуска http сервера")
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error(fmt.Sprintf("ошибка останова http сервера: %s", err))
	}
	<-ctx.Done()
	return nil
}

func (s *Server) HelloWorld(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, err := w.Write([]byte("hello-world")); err != nil {
		s.logger.Error(err.Error())
	}
}

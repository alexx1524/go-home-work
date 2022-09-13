package app

import (
	"context"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Storage interface {
	InsertEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID uuid.UUID) error
	GetEventByID(ctx context.Context, eventID uuid.UUID) (storage.Event, error)
	GetEventsCount(ctx context.Context) (int, error)
	GetEventsForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]storage.Event, error)
	RemoveEventsFinishedBeforeDate(ctx context.Context, date time.Time) (count int, err error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) InsertEvent(ctx context.Context, event storage.Event) error {
	return a.storage.InsertEvent(ctx, event)
}

func (a *App) UpdateEvent(ctx context.Context, event storage.Event) error {
	return a.storage.UpdateEvent(ctx, event)
}

func (a *App) DeleteEvent(ctx context.Context, eventID uuid.UUID) error {
	return a.storage.DeleteEvent(ctx, eventID)
}

func (a *App) GetEventByID(ctx context.Context, eventID uuid.UUID) (storage.Event, error) {
	return a.storage.GetEventByID(ctx, eventID)
}

func (a *App) GetEventsCount(ctx context.Context) (int, error) {
	return a.storage.GetEventsCount(ctx)
}

func (a *App) GetEventsForPeriod(ctx context.Context, dateStart time.Time, dateEnd time.Time) ([]storage.Event, error) {
	return a.storage.GetEventsForPeriod(ctx, dateStart, dateEnd)
}

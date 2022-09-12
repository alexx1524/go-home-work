package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Logger interface {
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

type Scheduler interface {
	Run(ctx context.Context) error
	Stop() error
	GetNotificationChannel() <-chan storage.Event
}

type scheduler struct {
	logger         Logger
	notificationCh chan storage.Event
	storage        app.Storage
}

func NewScheduler(logger Logger, appStorage app.Storage) Scheduler {
	return &scheduler{
		storage:        appStorage,
		logger:         logger,
		notificationCh: make(chan storage.Event),
	}
}

func (s *scheduler) Run(ctx context.Context) error {
	s.logger.Debug("Starting scheduler")

	go func(ctx context.Context) {
		s.logger.Info("Completed events removing is started")
		t := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Completed events removing is stopped")
				return
			case <-t.C:
				n, err := s.storage.RemoveEventsFinishedBeforeDate(ctx, time.Now().AddDate(-1, 0, 0))
				if err != nil {
					s.logger.Error(fmt.Sprintf("Removing completed events error %s", err.Error()))
				}
				s.logger.Debug(fmt.Sprintf("Removing completed events, removed %v events", n))
			}
		}
	}(ctx)

	go func(ctx context.Context) {
		s.logger.Info("Notifications watcher is started")
		t := time.NewTicker(1 * time.Hour)
		for {
			select {
			case <-ctx.Done():
				s.logger.Info("Notification watcher is stopped")
				return
			case <-t.C:
				now := time.Now()
				year, month, day := now.Date()
				begin := time.Date(year, month, day, 0, 0, 0, 0, now.Location())
				end := begin.AddDate(0, 0, 1)

				events, err := s.storage.GetEventsForPeriod(ctx, begin, end)
				if err != nil {
					s.logger.Error(fmt.Sprintf("Getting events error %s", err.Error()))
				}

				for _, e := range events {
					s.notificationCh <- e
				}
			}
		}
	}(ctx)

	return nil
}

func (s *scheduler) Stop() error {
	s.logger.Debug("Stopping scheduler")
	close(s.notificationCh)

	return nil
}

func (s *scheduler) GetNotificationChannel() <-chan storage.Event {
	return s.notificationCh
}

package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]storage.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]storage.Event),
	}
}

func (s *Storage) InsertEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; ok {
		return storage.ErrorEventAlreadyExists
	}

	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrorEventNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventByID(ctx context.Context, id uuid.UUID) (storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[id]

	if ok {
		return event, nil
	}

	return event, storage.ErrorEventNotFound
}

func (s *Storage) GetEventsCount(ctx context.Context) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.events), nil
}

func (s *Storage) GetEventsForPeriod(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := make([]storage.Event, 0)
	for _, value := range s.events {
		if value.StartDate.After(start) && value.StartDate.Before(end) {
			events = append(events, value)
		}
	}

	return events, nil
}

func (s *Storage) RemoveEventsFinishedBeforeDate(ctx context.Context, date time.Time) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	removeItems := make([]uuid.UUID, 0)
	for k, v := range s.events {
		if v.EndDate.Before(date) {
			removeItems = append(removeItems, k)
		}
	}

	for _, id := range removeItems {
		delete(s.events, id)
	}

	return len(removeItems), nil
}

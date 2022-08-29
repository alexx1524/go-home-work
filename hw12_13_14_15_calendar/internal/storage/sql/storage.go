package sqlstorage

import (
	"context"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	connectionString string
	db               *sqlx.DB
}

func New(connectionString string) (*Storage, error) {
	s := Storage{
		connectionString: connectionString,
	}

	err := s.Connect()
	return &s, err
}

func (s *Storage) Connect() error {
	db, err := sqlx.Open("postgres", s.connectionString)
	if err != nil {
		return err
	}

	s.db = db

	return err
}

func (s *Storage) Close() error {
	if err := s.db.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) InsertEvent(ctx context.Context, e storage.Event) error {
	query := `SELECT count(id) FROM events WHERE id = $1`

	var count int
	err := s.db.QueryRowContext(ctx, query, e.ID).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		return storage.ErrorEventAlreadyExists
	}

	query = `INSERT INTO events (id, title, description, user_id, start_date, end_date)
             VALUES($1, $2, $3, $4, $5, $6);`

	_, err = s.db.ExecContext(ctx, query, e.ID, e.Title, e.Description, e.UserID, e.StartDate, e.EndDate)

	return err
}

func (s *Storage) UpdateEvent(ctx context.Context, e storage.Event) error {
	query := `UPDATE events 
              SET title = $2, description = $3, user_id = $4, start_date = $5, end_date = $6
			  WHERE id = $1
              `

	_, err := s.db.ExecContext(ctx, query, e.ID, e.Title, e.Description, e.UserID, e.StartDate, e.EndDate)

	return err
}

func (s *Storage) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := s.db.ExecContext(ctx, query, id)

	return err
}

func (s *Storage) GetEventByID(ctx context.Context, id uuid.UUID) (event storage.Event, err error) {
	query := `SELECT id, title, description, user_id, start_date, end_date FROM events 
              WHERE id = $1`

	rows, err := s.db.QueryxContext(ctx, query, id)
	if err != nil {
		return event, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(&event)
		if err != nil {
			return event, err
		}
		return event, nil
	}

	return event, storage.ErrorEventNotFound
}

func (s *Storage) GetEventsCount(ctx context.Context) (count int, err error) {
	query := `SELECT count(id) FROM events`

	err = s.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (s *Storage) GetEventsForPeriod(ctx context.Context, start time.Time, end time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	query := `SELECT id, title, description, user_id, start_date, end_date 
              FROM events 
              WHERE start_date between $1 AND $2`

	rows, err := s.db.QueryxContext(ctx, query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		err = rows.StructScan(&event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

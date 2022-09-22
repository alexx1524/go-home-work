package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorEventNotFound      = errors.New("event not found")
	ErrorEventAlreadyExists = errors.New("event already exists")
	ErrorSourceUnsupported  = errors.New("source is unsupported")
)

type Event struct {
	ID          uuid.UUID `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	UserID      uuid.UUID `db:"user_id"`
	StartDate   time.Time `db:"start_date"`
	EndDate     time.Time `db:"end_date" `
}

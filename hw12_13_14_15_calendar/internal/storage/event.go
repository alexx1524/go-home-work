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
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	UserID      uuid.UUID `db:"user_id" json:"userId"`
	StartDate   time.Time `db:"start_date" json:"startDate"`
	EndDate     time.Time `db:"end_date" json:"endDate"`
}

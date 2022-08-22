package sqlstorage

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

func clearDB(ctx context.Context, connectionString string) error {
	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	db.ExecContext(ctx, "DELETE FROM events;")
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}

func TestStorage(t *testing.T) {
	connectionString := "postgres://pguser:pgpwd@localhost:5432/calendar?sslmode=disable"

	t.Run("insert new event", func(t *testing.T) {
		ctx := context.Background()
		err := clearDB(ctx, connectionString)
		if err != nil {
			log.Fatal(err)
		}

		sqlStorage, err := New(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlStorage.Close()

		event := storage.Event{
			ID:    uuid.New(),
			Title: "",
		}

		err = sqlStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		count, err := sqlStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("insert new event - already exists error", func(t *testing.T) {
		ctx := context.Background()
		err := clearDB(ctx, connectionString)
		if err != nil {
			log.Fatal(err)
		}

		sqlStorage, err := New(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlStorage.Close()

		event := storage.Event{
			ID:    uuid.New(),
			Title: "",
		}

		err = sqlStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		err = sqlStorage.InsertEvent(ctx, event)
		require.Error(t, err)
		require.True(t, errors.Is(err, storage.ErrorEventAlreadyExists))

		count, err := sqlStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("update event", func(t *testing.T) {
		ctx := context.Background()
		err := clearDB(ctx, connectionString)
		if err != nil {
			log.Fatal(err)
		}

		sqlStorage, err := New(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlStorage.Close()

		event := storage.Event{
			ID:    uuid.New(),
			Title: "title",
		}

		err = sqlStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		event.Title = "title1"
		err = sqlStorage.UpdateEvent(ctx, event)
		require.NoError(t, err)

		e, err := sqlStorage.GetEventByID(ctx, event.ID)
		require.NoError(t, err)
		require.Equal(t, event.Title, e.Title)

		count, err := sqlStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("delete event", func(t *testing.T) {
		ctx := context.Background()
		err := clearDB(ctx, connectionString)
		if err != nil {
			log.Fatal(err)
		}

		sqlStorage, err := New(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlStorage.Close()

		event := storage.Event{
			ID:    uuid.New(),
			Title: "title",
		}

		err = sqlStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		count, err := sqlStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)

		err = sqlStorage.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)

		count, err = sqlStorage.GetEventsCount(ctx)
		require.Equal(t, 0, count)
		require.NoError(t, err)
	})

	t.Run("get event by id - not found", func(t *testing.T) {
		ctx := context.Background()
		err := clearDB(ctx, connectionString)
		if err != nil {
			log.Fatal(err)
		}

		sqlStorage, err := New(connectionString)
		if err != nil {
			log.Fatal(err)
		}
		defer sqlStorage.Close()

		_, err = sqlStorage.GetEventByID(ctx, uuid.New())

		require.Error(t, err)
		require.True(t, errors.Is(err, storage.ErrorEventNotFound))
	})
}

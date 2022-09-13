package memorystorage

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alexx1524/go-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("insert new event", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		event := storage.Event{
			ID:    uuid.New(),
			Title: "",
		}

		err := memoryStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		count, err := memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("insert new event - already exists error", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		event := storage.Event{
			ID:    uuid.New(),
			Title: "",
		}

		err := memoryStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		err = memoryStorage.InsertEvent(ctx, event)
		require.Error(t, err)
		require.True(t, errors.Is(err, storage.ErrorEventAlreadyExists))

		count, err := memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("update event", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		event := storage.Event{
			ID:    uuid.New(),
			Title: "title",
		}

		err := memoryStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		event.Title = "title1"
		err = memoryStorage.UpdateEvent(ctx, event)
		require.NoError(t, err)

		e, err := memoryStorage.GetEventByID(ctx, event.ID)
		require.NoError(t, err)
		require.Equal(t, event.Title, e.Title)

		count, err := memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})

	t.Run("delete event", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		event := storage.Event{
			ID:    uuid.New(),
			Title: "title",
		}

		err := memoryStorage.InsertEvent(ctx, event)
		require.NoError(t, err)

		count, err := memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)

		err = memoryStorage.DeleteEvent(ctx, event.ID)
		require.NoError(t, err)

		count, err = memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 0, count)
		require.NoError(t, err)
	})

	t.Run("get for period", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		event := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Date(2022, time.August, 1, 1, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.September, 1, 1, 0, 0, 0, time.UTC),
		}
		event1 := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Date(2022, time.August, 1, 1, 10, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.September, 1, 1, 0, 0, 0, time.UTC),
		}
		event2 := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Date(2022, time.August, 2, 5, 10, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.September, 1, 1, 0, 0, 0, time.UTC),
		}
		event3 := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Date(2022, time.August, 5, 5, 10, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.September, 1, 1, 0, 0, 0, time.UTC),
		}

		err := memoryStorage.InsertEvent(ctx, event)
		require.NoError(t, err)
		err = memoryStorage.InsertEvent(ctx, event1)
		require.NoError(t, err)
		err = memoryStorage.InsertEvent(ctx, event2)
		require.NoError(t, err)
		err = memoryStorage.InsertEvent(ctx, event3)
		require.NoError(t, err)

		begin := time.Date(2022, time.August, 1, 0, 0, 0, 0, time.UTC)

		events, err := memoryStorage.GetEventsForPeriod(ctx, begin, begin.AddDate(0, 0, 1))
		require.NoError(t, err)
		require.Equal(t, 2, len(events))

		events, err = memoryStorage.GetEventsForPeriod(ctx, begin, begin.AddDate(0, 0, 7))
		require.NoError(t, err)
		require.Equal(t, 4, len(events))

		events, err = memoryStorage.GetEventsForPeriod(ctx, begin, begin.AddDate(0, 1, 0))
		require.NoError(t, err)
		require.Equal(t, 4, len(events))
	})
}

func TestRemoving(t *testing.T) {
	t.Run("remove completed events", func(t *testing.T) {
		memoryStorage := New()
		ctx := context.Background()
		oldEvent := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Date(2022, time.August, 1, 1, 0, 0, 0, time.UTC),
			EndDate:   time.Date(2022, time.September, 1, 1, 0, 0, 0, time.UTC),
		}
		event := storage.Event{
			ID:        uuid.New(),
			Title:     "title1",
			StartDate: time.Now(),
			EndDate:   time.Now().AddDate(0, 1, 0),
		}

		_ = memoryStorage.InsertEvent(ctx, oldEvent)
		_ = memoryStorage.InsertEvent(ctx, event)

		deletedCount, err := memoryStorage.RemoveEventsFinishedBeforeDate(ctx, time.Now().AddDate(0, 0, -1))

		require.Equal(t, 1, deletedCount)
		require.NoError(t, err)

		count, err := memoryStorage.GetEventsCount(ctx)
		require.Equal(t, 1, count)
		require.NoError(t, err)
	})
}

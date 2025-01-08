package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/helpers"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	s := New()

	e := storage.Event{
		Title:       "Some title",
		StartTime:   time.Now(),
		EndTime:     time.Now().Add(time.Hour),
		Description: helpers.StringToPointer("Some Description here"),
	}

	s.AddEvent(ctx, e)
	require.Equal(t, s.id, int64(1))
	require.NotEmpty(t, s.events[1])
}

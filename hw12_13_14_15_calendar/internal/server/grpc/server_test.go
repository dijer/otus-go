package grpcserver

import (
	"context"
	"testing"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/pb"
	factorystorage "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage/factory"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func setupTestServer() (*GRPCServer, error) {
	cfg := &config.Config{
		Logger: config.LoggerConf{
			Level: "INFO",
		},
		Storage: config.StorageConf{
			Storage: "memory",
		},
		GRPC: config.GRPCServerConf{
			Host:      "localhost",
			Port:      8080,
			Transport: "tcp",
		},
	}

	log := logger.New(cfg.Logger.Level)
	storage, err := factorystorage.New(cfg)
	if err != nil {
		return nil, err
	}

	app := app.New(log, storage)
	return New(log, app, cfg.GRPC), nil
}

func TestAddEvent(t *testing.T) {
	server, err := setupTestServer()
	require.NoError(t, err)

	req := &pb.AddEventRequest{
		Event: &pb.Event{
			Id:          1,
			Title:       "test title",
			Description: "test description",
			Owner:       123,
			StartTime:   timestamppb.New(time.Now()),
			EndTime:     timestamppb.New(time.Now().Add(time.Hour)),
		},
	}

	resp, err := server.AddEvent(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, "ok", resp.Message)
}

func TestUpdateEvent(t *testing.T) {
	server, err := setupTestServer()
	require.NoError(t, err)

	_, _ = server.AddEvent(context.Background(), &pb.AddEventRequest{
		Event: &pb.Event{
			Id:          1,
			Title:       "test title",
			Description: "test description",
			Owner:       123,
			StartTime:   timestamppb.New(time.Now()),
			EndTime:     timestamppb.New(time.Now().Add(time.Hour)),
		},
	})

	resp, err := server.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
		Event: &pb.Event{
			Id:          1,
			Title:       "updated title",
			Description: "updated description",
			Owner:       123,
			StartTime:   timestamppb.New(time.Now().Add(time.Hour)),
			EndTime:     timestamppb.New(time.Now().Add(2 * time.Hour)),
		},
	})

	require.NoError(t, err)
	require.Equal(t, "ok", resp.Message)
}

func TestDeleteEvent(t *testing.T) {
	server, err := setupTestServer()
	require.NoError(t, err)

	_, _ = server.AddEvent(context.Background(), &pb.AddEventRequest{
		Event: &pb.Event{
			Id:          1,
			Title:       "test title",
			Description: "test description",
			Owner:       123,
			StartTime:   timestamppb.New(time.Now()),
			EndTime:     timestamppb.New(time.Now().Add(time.Hour)),
		},
	})

	resp, err := server.DeleteEvent(context.Background(), &pb.DeleteEventRequest{
		Id: 1,
	})

	require.NoError(t, err)
	require.Equal(t, "ok", resp.Message)
}

func TestGetEventsList(t *testing.T) {
	server, err := setupTestServer()
	require.NoError(t, err)

	_, _ = server.AddEvent(context.Background(), &pb.AddEventRequest{
		Event: &pb.Event{
			Id:          1,
			Title:       "test title",
			Description: "test description",
			Owner:       123,
			StartTime:   timestamppb.New(time.Now()),
			EndTime:     timestamppb.New(time.Now().Add(time.Hour)),
		},
	})

	_, _ = server.AddEvent(context.Background(), &pb.AddEventRequest{
		Event: &pb.Event{
			Id:          2,
			Title:       "test title 2",
			Description: "test description 2",
			Owner:       321,
			StartTime:   timestamppb.New(time.Now().Add(time.Hour)),
			EndTime:     timestamppb.New(time.Now().Add(2 * time.Hour)),
		},
	})

	resp, err := server.GetEventsList(context.Background(), &pb.GetEventsListRequest{})
	require.NoError(t, err)
	require.Equal(t, "ok", resp.Message)
	require.Len(t, resp.Events, 2)
}

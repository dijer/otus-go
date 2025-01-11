package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/pb"
	grpcserver "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/stretchr/testify/require"
)

func getFreePort() (int, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	return listener.Addr().(*net.TCPAddr).Port, nil
}

func setupTestServers() (cancel context.CancelFunc, httpPort int, ctx context.Context) {
	httpPort, err := getFreePort()
	if err != nil {
		panic("failed get free port")
	}

	grpcPort, err := getFreePort()
	if err != nil {
		panic("failed get free port")
	}

	cfg := &config.Config{
		Logger: config.LoggerConf{
			Level: "INFO",
		},
		Storage: config.StorageConf{
			Storage: "memory",
		},
		GRPC: config.GRPCServerConf{
			Host:      "localhost",
			Port:      grpcPort,
			Transport: "tcp",
		},
		HTTP: config.HTTPServerConf{
			Host: "localhost",
			Port: httpPort,
		},
	}

	log := logger.New(cfg.Logger.Level)
	app := app.New(log, cfg)

	grpcSrv := grpcserver.New(log, app, cfg.GRPC)
	httpSrv := New(log, app, cfg.HTTP, cfg.GRPC)

	ctx, cancel = context.WithCancel(context.Background())

	go grpcSrv.Start(ctx)
	time.Sleep(500 * time.Millisecond)

	go httpSrv.Start(ctx)
	time.Sleep(500 * time.Millisecond)

	return
}

func TestHelloWorld(t *testing.T) {
	cancel, httpPort, ctx := setupTestServers()
	defer cancel()

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/hello",
		nil,
	)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestAddEvent(t *testing.T) {
	cancel, httpPort, _ := setupTestServers()
	defer cancel()

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/add",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "test title",
					"description": "test description",
					"start_time": "2025-01-11T10:00:00Z",
					"end_time": "2025-01-11T11:00:00Z",
					"owner": 1
				}
			}
		`)))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestUpdateEvent(t *testing.T) {
	cancel, httpPort, ctx := setupTestServers()
	defer cancel()

	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/add",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "test title",
					"description": "test description",
					"start_time": "2025-01-11T10:00:00Z",
					"end_time": "2025-01-11T11:00:00Z",
					"owner": 1
				}
			}
		`)))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/1",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "updated title",
					"description": "updated description",
					"start_time": "2025-01-11T12:00:00Z",
					"end_time": "2025-01-11T13:00:00Z",
					"owner": 1
				}
			}
		`)))
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestDeleteEvent(t *testing.T) {
	cancel, httpPort, ctx := setupTestServers()
	defer cancel()

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/add",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "test title",
					"description": "test description",
					"start_time": "2025-01-11T10:00:00Z",
					"end_time": "2025-01-11T11:00:00Z",
					"owner": 1
				}
			}
		`)))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/1",
		nil,
	)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetEventsList(t *testing.T) {
	cancel, httpPort, ctx := setupTestServers()
	defer cancel()

	req, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/add",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "test title",
					"description": "test description",
					"start_time": "2025-01-11T10:00:00Z",
					"end_time": "2025-01-11T11:00:00Z",
					"owner": 1
				}
			}
		`)))
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/add",
		bytes.NewBuffer([]byte(`
			{
				"event": {
					"title": "test title 2",
					"description": "test description 2",
					"start_time": "2025-01-11T11:00:00Z",
					"end_time": "2025-01-11T12:00:00Z",
					"owner": 123
				}
			}
		`)))
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, _ = http.NewRequestWithContext(
		ctx,
		http.MethodGet, "http://localhost:"+fmt.Sprintf("%d", httpPort)+"/events/list",
		nil,
	)
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var resBody pb.GetEventsListResponse
	err = json.NewDecoder(resp.Body).Decode(&resBody)
	require.NoError(t, err)
	require.Len(t, resBody.Events, 2)
}

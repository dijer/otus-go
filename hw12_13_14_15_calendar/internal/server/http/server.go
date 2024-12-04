package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	app    Application
	logger Logger
	config config.ServerConf
	server *http.Server
}

type Logger interface {
	Info(msg ...string)
	Error(msg ...string)
	Warn(msg ...string)
	Debug(msg ...string)
}

type Application interface {
	AddEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	GetEventsList(ctx context.Context) ([]storage.Event, error)
}

type Handler struct {
	logger Logger
}

func NewServer(logger Logger, app Application, config config.ServerConf) *Server {
	return &Server{
		logger: logger,
		app:    app,
		config: config,
	}
}

func (h *Handler) HelloWorld(w http.ResponseWriter, _ *http.Request) {
	h.logger.Info("Hello world!")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *Server) Start(_ context.Context) error {
	handler := &Handler{
		logger: s.logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", handler.loggingMiddleware(handler.HelloWorld))

	s.server = &http.Server{
		Addr:              net.JoinHostPort(s.config.Host, strconv.Itoa(int(s.config.Port))),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

package httpserver

import (
	"context"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/app"
	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/pb"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type HTTPServer struct {
	app     app.App
	logger  logger.Logger
	httpCfg config.HTTPServerConf
	grpcCfg config.GRPCServerConf
	server  *http.Server
}

type Handler struct {
	logger logger.Logger
}

const readHeaderTimeout = 5 * time.Second

func New(logger *logger.Logger, app app.App, httpCfg config.HTTPServerConf, grpcCfg config.GRPCServerConf) *HTTPServer {
	return &HTTPServer{
		logger:  *logger,
		app:     app,
		httpCfg: httpCfg,
		grpcCfg: grpcCfg,
	}
}

func (h *Handler) HelloWorld(w http.ResponseWriter, _ *http.Request, _ map[string]string) {
	h.logger.Info("Hello world!")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}

func (s *HTTPServer) Start(ctx context.Context) error {
	handler := &Handler{
		logger: s.logger,
	}

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	)
	mux.HandlePath("GET", "/hello", handler.loggingMiddleware(handler.HelloWorld))

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterCalendarServiceHandlerFromEndpoint(
		ctx,
		mux,
		net.JoinHostPort(s.grpcCfg.Host, strconv.Itoa(s.grpcCfg.Port)),
		opts,
	)
	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:              net.JoinHostPort(s.httpCfg.Host, strconv.Itoa(s.httpCfg.Port)),
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	log.Println("HTTP server started on", s.httpCfg.Host, ":", s.httpCfg.Port)
	return s.server.ListenAndServe()
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

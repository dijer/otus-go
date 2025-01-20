package grpcserver

import (
	"context"
	"net"
	"strconv"

	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/app"
	config "github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/config/calendar"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/pb"
	"github.com/dijer/otus-go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCServer struct {
	pb.UnimplementedCalendarServiceServer
	server *grpc.Server
	app    app.App
	log    logger.Logger
	config config.GRPCServerConf
}

func New(
	logger *logger.Logger,
	app app.App,
	cfg config.GRPCServerConf,
) *GRPCServer {
	grpcServer := grpc.NewServer()

	grpcServerStruct := &GRPCServer{
		server: grpcServer,
		app:    app,
		log:    *logger,
		config: cfg,
	}

	pb.RegisterCalendarServiceServer(grpcServer, grpcServerStruct)

	return grpcServerStruct
}

func (s *GRPCServer) Start(_ context.Context) error {
	lsn, err := net.Listen(s.config.Transport, net.JoinHostPort(s.config.Host, strconv.Itoa(s.config.Port)))
	if err != nil {
		return err
	}

	s.log.Info("strt grpc")

	return s.server.Serve(lsn)
}

func (s *GRPCServer) Stop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}

func (s *GRPCServer) AddEvent(ctx context.Context, req *pb.AddEventRequest) (*pb.AddEventResponse, error) {
	if req == nil || req.Event == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event missing")
	}

	event := storage.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		Owner:       req.Event.Owner,
		StartTime:   req.Event.StartTime.AsTime(),
		EndTime:     req.Event.EndTime.AsTime(),
		Description: &req.Event.Description,
	}
	err := s.app.AddEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	resp := &pb.AddEventResponse{
		Message: "ok",
	}

	return resp, nil
}

func (s *GRPCServer) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.UpdateEventResponse, error) {
	event := storage.Event{
		ID:          req.Event.Id,
		Title:       req.Event.Title,
		StartTime:   req.Event.StartTime.AsTime(),
		EndTime:     req.Event.EndTime.AsTime(),
		Description: &req.Event.Description,
		Owner:       req.Event.Owner,
	}

	err := s.app.UpdateEvent(ctx, event)
	if err != nil {
		return nil, err
	}

	resp := &pb.UpdateEventResponse{
		Message: "ok",
	}

	return resp, nil
}

func (s *GRPCServer) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	id := req.Id
	err := s.app.DeleteEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	resp := &pb.DeleteEventResponse{
		Message: "ok",
	}

	return resp, nil
}

func (s *GRPCServer) GetEventsList(ctx context.Context, _ *pb.GetEventsListRequest) (*pb.GetEventsListResponse, error) {
	events, err := s.app.GetEventsList(ctx)
	if err != nil {
		return nil, err
	}

	pbEvents := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Id:          event.ID,
			Title:       event.Title,
			Description: *event.Description,
			StartTime:   timestamppb.New(event.StartTime),
			EndTime:     timestamppb.New(event.EndTime),
			Owner:       event.Owner,
		})
	}

	resp := &pb.GetEventsListResponse{
		Events:  pbEvents,
		Message: "ok",
	}

	return resp, nil
}

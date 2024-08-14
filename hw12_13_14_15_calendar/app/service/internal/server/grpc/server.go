package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/common"
	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/logger"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	host  string
	port  string
	store storage.Storage
	pb.UnimplementedEventServiceServer
}

func New(app *app.App) *GrpcServer {
	return &GrpcServer{
		host:  Config.GrpcServer.Host,
		port:  strconv.Itoa(Config.GrpcServer.Port),
		store: app.Storage(),
	}
}

func (s *GrpcServer) Start() error {
	lis, err := net.Listen("tcp", `:`+s.port)
	if err != nil {
		return err
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(newLoggingInterceptor()))

	pb.RegisterEventServiceServer(srv, s)

	logger.Log.Info("Starting gRPC server at " + s.host + ":" + s.port)

	if err := srv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func (s *GrpcServer) Stop(_ context.Context) error {
	return nil
}

func (s *GrpcServer) AddEvent(_ context.Context, req *pb.AddEventRequest) (*pb.Event, error) {
	if req.Event == nil {
		return nil, fmt.Errorf("there is no event in request")
	}
	eventProto := EventProto{
		Title: req.Event.Title,
		Time:  req.Event.Time,
	}
	event := eventProto.FromProto()
	err := s.store.AddEvent(*event)
	if err != nil {
		return nil, err
	}
	return req.Event, nil
}

func (s *GrpcServer) UpdateEvent(_ context.Context, req *pb.UpdateEventRequest) (*pb.Event, error) {
	if req.Event == nil {
		return nil, fmt.Errorf("there is no event in request")
	}
	eventProto := EventProto{
		ID:    req.Event.Id,
		Title: req.Event.Title,
		Time:  req.Event.Time,
	}
	event := eventProto.FromProto()
	err := s.store.UpdateEvent(*event)
	if err != nil {
		return nil, err
	}
	return req.Event, nil
}

func (s *GrpcServer) DeleteEvent(_ context.Context, req *pb.DeleteEventRequest) (*pb.Event, error) {
	err := s.store.DeleteEvent(int(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.Event{Id: req.Id}, nil
}

func (s *GrpcServer) ListEvents(_ context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, status.Errorf(http.StatusBadRequest, "invalid start date: %v", err)
	}

	events, err := s.store.ListEvents(startDate, int(req.Days))
	if err != nil {
		return nil, err
	}
	var pbEvents []*pb.Event
	for _, event := range events {
		pbEvents = append(pbEvents, &pb.Event{
			Id:               int32(event.ID),
			Title:            event.Title,
			Time:             timestamppb.New(event.Time),
			NotificationSent: event.NotificationSent,
		})
	}
	return &pb.ListEventsResponse{Events: pbEvents}, nil
}

func newLoggingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		logger.Log.Infof("gRPC Request received: method=%s, request=%v", info.FullMethod, req)

		start := time.Now()
		resp, err = handler(ctx, req)
		duration := time.Since(start)

		if err != nil {
			logger.Log.Errorf("gRPC Request failed: method=%s, request=%v, error=%v, duration=%s, status=%s",
				info.FullMethod, req, err, duration, status.Code(err))
		} else {
			logger.Log.Infof("gRPC Request succeeded: method=%s, request=%v, response=%v, duration=%s",
				info.FullMethod, req, resp, duration)
		}
		return resp, err
	}
}

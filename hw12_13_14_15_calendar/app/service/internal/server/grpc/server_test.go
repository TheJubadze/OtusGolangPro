package internalgrpc

import (
	"context"
	"net"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/proto/pb"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGrpcServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage.NewMockStorage(ctrl)
	application := app.New(mockStorage)
	server := New(application)

	s := grpc.NewServer()
	pb.RegisterEventServiceServer(s, server)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Errorf("Server exited with error: %v", err)
			return
		}
	}()
	defer s.Stop()

	conn, err := grpc.Dial("bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials())) // nolint:staticcheck
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	client := pb.NewEventServiceClient(conn)

	t.Run("AddEvent", func(t *testing.T) {
		mockStorage.EXPECT().AddEvent(gomock.Any()).Return(nil).Times(1)

		_, err := client.AddEvent(context.Background(), &pb.AddEventRequest{
			Event: &pb.Event{Title: "Test Event", Time: timestamppb.New(time.Now())},
		})
		assert.NoError(t, err)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		mockStorage.EXPECT().UpdateEvent(gomock.Any()).Return(nil).Times(1)

		_, err := client.UpdateEvent(context.Background(), &pb.UpdateEventRequest{
			Event: &pb.Event{Id: 1, Title: "Updated Event", Time: timestamppb.New(time.Now())},
		})
		assert.NoError(t, err)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		mockStorage.EXPECT().DeleteEvent(gomock.Any()).Return(nil).Times(1)

		_, err := client.DeleteEvent(context.Background(), &pb.DeleteEventRequest{Id: 1})
		assert.NoError(t, err)
	})

	t.Run("ListEvents", func(t *testing.T) {
		mockStorage.EXPECT().ListEvents(gomock.Any(), gomock.Any()).Return([]Event{
			{ID: 1, Title: "Event 1", Time: time.Now()},
			{ID: 2, Title: "Event 2", Time: time.Now()},
		}, nil).Times(1)

		resp, err := client.ListEvents(context.Background(), &pb.ListEventsRequest{StartDate: "2023-01-01", Days: 7})
		assert.NoError(t, err)
		assert.Len(t, resp.Events, 2)
	})
}

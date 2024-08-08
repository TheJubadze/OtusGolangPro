package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/proto/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	address = "localhost:8080"
)

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)
	client := pb.NewEventServiceClient(conn)

	for {
		fmt.Println("Choose an action: [1] Add Event [2] Update Event [3] Delete Event [4] List Events [5] Exit")
		var choice int
		_, _ = fmt.Scan(&choice)

		switch choice {
		case 1:
			addEvent(client)
		case 2:
			updateEvent(client)
		case 3:
			deleteEvent(client)
		case 4:
			listEvents(client)
		case 5:
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func addEvent(client pb.EventServiceClient) {
	var title string
	var timestamp int64

	fmt.Print("Enter event title: ")
	_, _ = fmt.Scan(&title)
	fmt.Print("Enter event time (Unix timestamp): ")
	_, _ = fmt.Scan(&timestamp)

	event := &pb.Event{
		Title: title,
		Time:  timestamppb.New(time.Unix(timestamp, 0)),
	}

	_, err := client.AddEvent(context.Background(), &pb.AddEventRequest{Event: event})
	if err != nil {
		log.Printf("could not add event: %v\n", err)
		return
	}
	fmt.Println("Event added successfully!")
}

func updateEvent(client pb.EventServiceClient) {
	var id int32
	var title string
	var timestamp int64

	fmt.Print("Enter event ID: ")
	_, _ = fmt.Scan(&id)
	fmt.Print("Enter event title: ")
	_, _ = fmt.Scan(&title)
	fmt.Print("Enter event time (Unix timestamp): ")
	_, _ = fmt.Scan(&timestamp)

	event := &pb.Event{
		Id:    id,
		Title: title,
		Time:  timestamppb.New(time.Unix(timestamp, 0)),
	}

	_, err := client.UpdateEvent(context.Background(), &pb.UpdateEventRequest{Event: event})
	if err != nil {
		log.Printf("could not update event: %v\n", err)
		return
	}
	fmt.Println("Event updated successfully!")
}

func deleteEvent(client pb.EventServiceClient) {
	var id int32

	fmt.Print("Enter event ID: ")
	_, _ = fmt.Scan(&id)

	_, err := client.DeleteEvent(context.Background(), &pb.DeleteEventRequest{Id: id})
	if err != nil {
		log.Printf("could not delete event: %v\n", err)
		return
	}
	fmt.Println("Event deleted successfully!")
}

func listEvents(client pb.EventServiceClient) {
	var startDate string
	var days int32

	fmt.Print("Enter start date: ")
	_, _ = fmt.Scan(&startDate)
	fmt.Print("Enter days count: ")
	_, _ = fmt.Scan(&days)

	resp, err := client.ListEvents(context.Background(), &pb.ListEventsRequest{
		StartDate: startDate,
		Days:      days,
	})
	if err != nil {
		log.Printf("could not list events: %v\n", err)
		return
	}

	fmt.Println("Events:")
	for _, event := range resp.Events {
		fmt.Printf("ID: %d, Title: %s, Time: %s\n", event.Id, event.Title, event.Time.AsTime().Format(time.RFC3339))
	}
}

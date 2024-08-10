package entity

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type Event struct {
	ID               int       `db:"id"`
	Title            string    `db:"title"`
	Time             time.Time `db:"time"`
	NotificationSent bool      `db:"notification_sent"`
}

// EventProto is the protobuf representation of Event
type EventProto struct {
	ID    int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Title string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Time  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=time,proto3" json:"time,omitempty"`
}

// ToProto converts Event to its protobuf representation
func (e *Event) ToProto() *EventProto {
	return &EventProto{
		ID:    int32(e.ID),
		Title: e.Title,
		Time:  timestamppb.New(e.Time),
	}
}

// FromProto converts EventProto to Event
func (proto *EventProto) FromProto() *Event {
	e := Event{}
	e.ID = int(proto.ID)
	e.Title = proto.Title
	e.Time = proto.Time.AsTime()
	return &e
}

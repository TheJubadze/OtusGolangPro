package storage

import (
	"time"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"
)

type Storage interface {
	AddEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(id int) error
	ListEvents(time.Time, int) ([]Event, error)
	ListEventsToNotify() ([]Event, error)
	MarkEventsNotificationSent(ids []int) error
	DeleteEventsOlderThan(years, months, days int) error
}

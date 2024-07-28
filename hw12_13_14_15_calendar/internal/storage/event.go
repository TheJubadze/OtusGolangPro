package storage

type Event struct {
	ID    int
	Title string
	Time  string
}

type Storage interface {
	AddEvent(event Event) error
	UpdateEvent(event Event) error
	DeleteEvent(id int) error
	ListEvents() ([]Event, error)
}

package memorystorage

import (
	"errors"
	"sync"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/storage"
)

type Event = storage.Event

type InMemoryStorage struct {
	mu     sync.Mutex
	events map[int]Event
}

func New() *InMemoryStorage {
	return &InMemoryStorage{
		events: make(map[int]Event),
	}
}

func (s *InMemoryStorage) AddEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.ID]; exists {
		return errors.New("event already exists")
	}
	s.events[event.ID] = event
	return nil
}

func (s *InMemoryStorage) UpdateEvent(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[event.ID]; !exists {
		return errors.New("event not found")
	}
	s.events[event.ID] = event
	return nil
}

func (s *InMemoryStorage) DeleteEvent(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.events[id]; !exists {
		return errors.New("event not found")
	}
	delete(s.events, id)
	return nil
}

func (s *InMemoryStorage) ListEvents() ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	return events, nil
}

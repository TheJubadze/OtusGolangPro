package memorystorage

import (
	"errors"
	"sync"
	"time"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"
)

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

func (s *InMemoryStorage) ListEvents(startDate time.Time, days int) ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]Event, 0, len(s.events))
	for _, event := range s.events {
		if event.Time.Before(startDate) || event.Time.After(startDate.AddDate(0, 0, days)) {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *InMemoryStorage) ListEventsToNotify() ([]Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	events := make([]Event, 0, len(s.events))
	for _, event := range s.events {
		if event.NotificationSent {
			continue
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *InMemoryStorage) MarkEventsNotificationSent(ids []int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		event, exists := s.events[id]
		if !exists {
			return errors.New("event not found")
		}
		event.NotificationSent = true
		s.events[id] = event
	}
	return nil
}

func (s *InMemoryStorage) DeleteEventsOlderThan(years, months, days int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, event := range s.events {
		if event.Time.Before(time.Now().AddDate(-years, -months, -days)) {
			delete(s.events, id)
		}
	}
	return nil
}

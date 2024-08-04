package memorystorage

import (
	"math"
	"testing"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var timeNow = time.Now()

func TestAddEvent(t *testing.T) {
	s := New()
	event := storage.Event{ID: 1, Title: "Test Event", Time: timeNow}

	// Test adding a new event
	err := s.AddEvent(event)
	require.NoError(t, err)

	// Test adding a duplicate event
	err = s.AddEvent(event)
	require.Error(t, err)
	require.Equal(t, "event already exists", err.Error())
}

func TestUpdateEvent(t *testing.T) {
	s := New()
	event := storage.Event{ID: 1, Title: "Test Event", Time: timeNow}

	// Test updating a non-existent event
	nonExistentEvent := storage.Event{ID: 2, Title: "Non-existent Event", Time: timeNow}
	err := s.UpdateEvent(nonExistentEvent)
	require.Error(t, err)
	require.Equal(t, "event not found", err.Error())

	// Test updating an existing event
	err = s.AddEvent(event)
	require.NoError(t, err)
	updatedEvent := storage.Event{ID: 1, Title: "Updated Event", Time: timeNow}
	err = s.UpdateEvent(updatedEvent)
	require.NoError(t, err)

	// Verify update
	storedEvent := s.events[event.ID]
	require.Equal(t, "Updated Event", storedEvent.Title)
	require.Equal(t, timeNow, storedEvent.Time)
}

func TestDeleteEvent(t *testing.T) {
	s := New()
	event := storage.Event{ID: 1, Title: "Test Event", Time: timeNow}

	// Test deleting a non-existent event
	err := s.DeleteEvent(event.ID)
	require.Error(t, err)
	require.Equal(t, "event not found", err.Error())

	// Test deleting an existing event
	err = s.AddEvent(event)
	require.NoError(t, err)
	err = s.DeleteEvent(event.ID)
	require.NoError(t, err)

	// Verify deletion
	_, exists := s.events[event.ID]
	require.False(t, exists)
}

func TestListEvents(t *testing.T) {
	s := New()
	event1 := storage.Event{ID: 1, Title: "Test Event 1", Time: timeNow}
	event2 := storage.Event{ID: 2, Title: "Test Event 2", Time: timeNow}

	// Test listing events in an empty storage
	events, err := s.ListEvents(time.Time{}, math.MaxInt32)
	require.NoError(t, err)
	require.Len(t, events, 0)

	// Test listing events after adding events
	err = s.AddEvent(event1)
	require.NoError(t, err)
	err = s.AddEvent(event2)
	require.NoError(t, err)
	events, err = s.ListEvents(time.Time{}, math.MaxInt32)
	require.NoError(t, err)
	require.Len(t, events, 2)

	// Verify the events are correctly listed
	require.Equal(t, event1.ID, events[0].ID)
	require.Equal(t, event2.ID, events[1].ID)
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	t1 := time.Date(2024, time.July, 28, 16, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, time.July, 28, 16, 0, 0, 0, time.UTC)
	event1 := storage.Event{ID: 1, Title: "Test Event 1", Time: t1}
	event2 := storage.Event{ID: 2, Title: "Test Event 2", Time: t2}

	done := make(chan bool)

	// Test concurrent add events
	go func() {
		err := s.AddEvent(event1)
		require.NoError(t, err)
		done <- true
	}()
	go func() {
		err := s.AddEvent(event2)
		require.NoError(t, err)
		done <- true
	}()
	<-done
	<-done

	// Verify both events were added
	events, err := s.ListEvents(time.Time{}, math.MaxInt32)
	require.NoError(t, err)
	require.Len(t, events, 2)

	// Test concurrent update and delete events
	go func() {
		updatedEvent := storage.Event{ID: 1, Title: "Updated Event 1", Time: t1}
		err := s.UpdateEvent(updatedEvent)
		require.NoError(t, err)
		done <- true
	}()
	go func() {
		err := s.DeleteEvent(2)
		require.NoError(t, err)
		done <- true
	}()
	<-done
	<-done

	// Verify update and delete
	events, err = s.ListEvents(time.Time{}, math.MaxInt32)
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "Updated Event 1", events[0].Title)
	require.Equal(t, t1, events[0].Time)
}

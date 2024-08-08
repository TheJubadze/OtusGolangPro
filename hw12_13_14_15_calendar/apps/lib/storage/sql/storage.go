package sqlstorage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/lib/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pressly/goose"
)

type SQLStorage struct {
	db *sqlx.DB
}

func New(dsn string) (*SQLStorage, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &SQLStorage{db: db}, nil
}

func (s *SQLStorage) Migrate(migrationsDir string) (err error) {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}
	if err := goose.Up(s.db.DB, migrationsDir); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *SQLStorage) Close() error {
	return s.db.Close()
}

func (s *SQLStorage) AddEvent(event storage.Event) error {
	_, err := s.db.NamedExecContext(context.Background(), "INSERT INTO events (title, time) VALUES (:title, :time)", event)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLStorage) UpdateEvent(event storage.Event) error {
	result, err := s.db.NamedExecContext(context.Background(), "UPDATE events SET title = :title, time = :time WHERE id = :id", event)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (s *SQLStorage) DeleteEvent(id int) error {
	result, err := s.db.ExecContext(context.Background(), "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (s *SQLStorage) ListEvents(startDate time.Time, daysCount int) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.SelectContext(context.Background(), &events, "SELECT id, title, time FROM events where time between $1 and $2", startDate, startDate.AddDate(0, 0, daysCount))
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) ListEventsToNotify() ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.SelectContext(context.Background(), &events, "SELECT id, title, time FROM events where time < $1 and not notification_sent", time.Now())
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *SQLStorage) MarkEventsNotificationSent(ids []int) error {
	result, err := s.db.ExecContext(context.Background(), "UPDATE events SET notification_sent = true WHERE id IN $1", ids)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

func (s *SQLStorage) DeleteEventsOlderThan(years, months, days int) error {
	result, err := s.db.ExecContext(context.Background(), "DELETE FROM events WHERE time > $1", time.Now().AddDate(-years, -months, -days))
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}
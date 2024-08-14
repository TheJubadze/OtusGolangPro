package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/common"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/amqp"
	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage/sql"
)

var (
	configFile string
	cfg        = Config
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if err := Init(configFile); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	sqlStorage, err := sqlstorage.New(cfg.Storage.DSN)
	if err != nil {
		log.Fatalf("Error initializing storage: %s", err)
	}

	publisher := amqp.NewPublisher(cfg.Amqp)

	ticker := time.NewTicker(cfg.Amqp.NotificationPeriodSeconds * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			notify(sqlStorage, publisher)
		}
	}()

	select {}
}

func notify(storage Storage, publisher *amqp.Publisher) {
	eventsToNotify, err := storage.ListEventsToNotify()
	if err != nil {
		log.Fatalf("Error getting events to notify: %s", err)
	}

	for _, event := range eventsToNotify {
		eventJSON, err := json.Marshal(struct {
			ID    int       `json:"id"`
			Title string    `json:"title"`
			Time  time.Time `json:"time"`
		}{
			ID:    event.ID,
			Title: event.Title,
			Time:  event.Time,
		})
		if err != nil {
			log.Fatalf("Error marshaling event: %s", err)
		}
		if err := publisher.Publish(string(eventJSON)); err != nil {
			log.Fatalf("%s", err)
		}
	}

	ids := make([]int, len(eventsToNotify))
	for i, event := range eventsToNotify {
		ids[i] = event.ID
	}

	err = storage.MarkEventsNotificationSent(ids)
	if err != nil {
		log.Fatalf("Error marking events as notified: %s", err)
	}

	err = storage.DeleteEventsOlderThan(1, 0, 0)
	if err != nil {
		log.Fatalf("Error deleting old events: %s", err)
	}
}

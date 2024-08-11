package main

import (
	"flag"
	"log"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/common"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/amqp"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage/sql"
)

var (
	configFile string
	cfg        = Config
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/scheduler/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if err := Init(configFile); err != nil {
		log.Fatalf("Error reading config: %s", err)
	}

	_, err := sqlstorage.New(cfg.Storage.DSN)
	if err != nil {
		log.Fatalf("Error initializing storage: %s", err)
	}

	consumer := amqp.NewConsumer(cfg.Amqp)

	err = consumer.Consume()
	if err != nil {
		log.Fatalf("Error consuming: %s", err)
	}
}

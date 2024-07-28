package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/config"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/errors"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/storage/sql"
)

var (
	configFile string
	cfg        = config.Config
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if err := config.Init(configFile); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		os.Exit(1)
	}

	logger.Setup(cfg.Logger.Level)

	storage, err := initStorage()
	if err != nil {
		logger.Log().Fatalf("Error initializing storage: %s", err)
	}

	calendar := app.New(storage)
	server, err := internalhttp.New(calendar)
	if err != nil {
		logger.Log().Fatalf("Error initializing http server: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Log().Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Log().Info("calendar is running...")

	if err := server.Start(); err != nil {
		logger.Log().Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func initStorage() (app.Storage, error) {
	var storage app.Storage
	switch cfg.Storage.Type {
	case "memory":
		storage = memorystorage.New()
	case "sql":
		sqlStorage, err := sqlstorage.New(cfg.Storage.DSN)
		if err != nil {
			return nil, &errors.ProgramError{Err: fmt.Errorf("error initializing sql storage: %w", err)}
		}

		err = sqlStorage.Migrate(cfg.Storage.MigrationsDir)
		if err != nil {
			return nil, &errors.ProgramError{Err: fmt.Errorf("error migrating sql storage: %w", err)}
		}

		storage = sqlStorage

		defer func() {
			if err := sqlStorage.Close(); err != nil {
				log.Println("cannot close psql connection", err)
			}
		}()
	default:
		return nil, &errors.ProgramError{Err: fmt.Errorf("error reading storage config: unknown storage type '%s'", cfg.Storage.Type)}
	}

	return storage, nil
}

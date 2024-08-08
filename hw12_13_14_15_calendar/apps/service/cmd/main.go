package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/common"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/lib/storage/memory"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/lib/storage/sql"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/config"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/logger"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/server"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/server/grpc"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/apps/service/internal/server/http"
)

var (
	configFile = flag.String("config", "/etc/calendar/config.yaml", "Path to configuration file")
	cfg        = config.Config
)

func init() {
	flag.Parse()
}

func main() {
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if err := config.Init(*configFile); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error reading config: %s\n", err)
		os.Exit(1)
	}

	logger.Setup(cfg.Logger.Level)

	storage, closeFunc, err := initStorage()
	defer closeFunc()
	if err != nil {
		logger.Log.Fatalf("Error initializing storage: %s", err)
	}

	calendar := app.New(storage)
	servers := []serverinterface.Server{internalhttp.New(calendar), internalgrpc.New(calendar)}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	for _, server := range servers {
		wg.Add(1)
		go startServer(server, cancel, &wg)
	}

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		for _, server := range servers {
			if err := server.Stop(ctx); err != nil {
				logger.Log.Error("failed to stop server: " + err.Error())
			}
		}
	}()

	wg.Wait()
}

func startServer(server serverinterface.Server, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := server.Start(); err != nil {
		logger.Log.Error("failed to start server: " + err.Error())
		cancel()
	}
}

func initStorage() (app.Storage, func(), error) {
	var storage app.Storage
	var closeFunc func()
	switch cfg.Storage.Type {
	case "memory":
		storage = memorystorage.New()
		closeFunc = func() {}
	case "sql":
		sqlStorage, err := sqlstorage.New(cfg.Storage.DSN)
		if err != nil {
			return nil, nil, &common.ProgramError{Err: fmt.Errorf("error initializing sql storage: %w", err)}
		}

		logger.Log.Println("connected to psql")

		err = sqlStorage.Migrate(cfg.Storage.MigrationsDir)
		if err != nil {
			return nil, nil, &common.ProgramError{Err: fmt.Errorf("error migrating sql storage: %w", err)}
		}

		logger.Log.Println("migrated psql")

		storage = sqlStorage

		closeFunc = func() {
			if err := sqlStorage.Close(); err != nil {
				logger.Log.Println("cannot close psql connection", err)
			}
		}
	default:
		return nil, nil, &common.ProgramError{Err: fmt.Errorf("error reading storage config: unknown storage type '%s'", cfg.Storage.Type)}
	}

	return storage, closeFunc, nil
}

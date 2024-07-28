package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/config"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/logger"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/internal/storage"
)

var cfg = config.Config

type Server struct {
	host   string
	port   int
	store  storage.Storage
	server *http.Server
}

func New(app *app.App) (*Server, error) {
	s := &Server{
		host:  cfg.Server.Host,
		port:  cfg.Server.Port,
		store: app.Storage(),
	}

	loggedMux := LoggingMiddleware(s.routes())

	s.server = &http.Server{
		Addr:              fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:           loggedMux,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
		ReadTimeout:       15 * time.Second,
	}
	return s, nil
}

func (s *Server) Start() error {
	logger.Log().Printf("Starting server at %s:%d", s.host, s.port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	logger.Log().Println("Shutting down server...")
	return s.server.Shutdown(ctx)
}

func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	return mux
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		return
	}
}

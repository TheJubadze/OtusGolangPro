package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/common"
	. "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/app"
	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/logger"
)

type HttpServer struct {
	host   string
	port   int
	store  storage.Storage
	server *http.Server
}

func New(app *app.App) *HttpServer {
	s := &HttpServer{
		host:  Config.HttpServer.Host,
		port:  Config.HttpServer.Port,
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
	return s
}

func (s *HttpServer) Start() error {
	logger.Log.Printf("Starting HTTP server at %s:%d", s.host, s.port)
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	logger.Log.Println("Shutting down server...")
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /hello", helloHandler)
	addEventsRoutes(mux, s)
	return mux
}

func addEventsRoutes(mux *http.ServeMux, s *HttpServer) {
	mux.HandleFunc("POST /events", s.addEventHandler)
	mux.HandleFunc("PUT /events", s.updateEventHandler)
	mux.HandleFunc("DELETE /events", s.deleteEventHandler)
	mux.HandleFunc("GET /events", s.listEventsHandler)
}

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Hello, World!"))
	if err != nil {
		return
	}
}

func (s *HttpServer) addEventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newEvent := Event{
		Title: event.Title,
		Time:  event.Time,
	}

	if err := s.store.AddEvent(newEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HttpServer) updateEventHandler(w http.ResponseWriter, r *http.Request) {
	var event Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEvent := Event{
		ID:    event.ID,
		Title: event.Title,
		Time:  event.Time,
	}

	if err := s.store.UpdateEvent(updatedEvent); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HttpServer) deleteEventHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := s.store.DeleteEvent(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := Event{ID: id}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *HttpServer) listEventsHandler(w http.ResponseWriter, r *http.Request) {
	startDate, err := time.Parse("2006-01-02", r.URL.Query().Get("startDate"))
	if err != nil {
		http.Error(w, "Invalid start date", http.StatusBadRequest)
		return
	}
	daysCount, err := strconv.Atoi(r.URL.Query().Get("days"))
	if err != nil {
		http.Error(w, "Invalid days count", http.StatusBadRequest)
		return
	}
	events, err := s.store.ListEvents(startDate, daysCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var eventList []Event
	for _, event := range events {
		eventList = append(eventList, Event{
			ID:               event.ID,
			Title:            event.Title,
			Time:             event.Time,
			NotificationSent: event.NotificationSent,
		})
	}

	if err := json.NewEncoder(w).Encode(eventList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

package internalhttp

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/entity"

	storage2 "github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/lib/storage"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/app"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHttpServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := storage2.NewMockStorage(ctrl)
	application := app.New(mockStorage)
	server := New(application)

	t.Run("addEventHandler", func(t *testing.T) {
		event := entity.Event{Title: "Test Event", Time: time.Now()}
		mockStorage.EXPECT().AddEvent(gomock.Any()).Return(nil).Times(1)

		body, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.addEventHandler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("updateEventHandler", func(t *testing.T) {
		event := entity.Event{ID: 1, Title: "Updated Event", Time: time.Now()}
		mockStorage.EXPECT().UpdateEvent(gomock.Any()).Return(nil).Times(1)

		body, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPut, "/events", bytes.NewReader(body))
		w := httptest.NewRecorder()

		server.updateEventHandler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("deleteEventHandler", func(t *testing.T) {
		id := 1
		mockStorage.EXPECT().DeleteEvent(id).Return(nil).Times(1)

		req := httptest.NewRequest(http.MethodDelete, "/events?id="+strconv.Itoa(id), nil)
		w := httptest.NewRecorder()

		server.deleteEventHandler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("listEventsHandler", func(t *testing.T) {
		startDate := time.Now()
		days := 7
		events := []entity.Event{
			{ID: 1, Title: "Event 1", Time: startDate},
			{ID: 2, Title: "Event 2", Time: startDate.Add(24 * time.Hour)},
		}
		mockStorage.EXPECT().ListEvents(gomock.Any(), days).Return(events, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/events?startDate="+startDate.Format("2006-01-02")+"&days="+strconv.Itoa(days), nil)
		w := httptest.NewRecorder()

		server.listEventsHandler(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

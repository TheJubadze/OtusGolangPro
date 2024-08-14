package tests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

type Event struct {
	ID               int       `db:"id"`
	Title            string    `db:"title"`
	Time             time.Time `db:"time"`
	NotificationSent bool      `db:"notification_sent"`
}

var (
	client  *http.Client
	baseURL string
)

var _ = ginkgo.BeforeSuite(func() {
	client = &http.Client{}
	baseURL = "http://calendar:8080"
})

var _ = ginkgo.Describe("Calendar API Integration Tests", func() {
	ginkgo.Context("When interacting with the calendar API", func() {
		ginkgo.It("should create an event successfully", func() {
			event := Event{Title: "Test Event", Time: time.Now()}
			body, _ := json.Marshal(event)
			req, err := http.NewRequest("POST", baseURL+"/events", bytes.NewReader(body))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusCreated))
		})

		ginkgo.It("should update an event successfully", func() {
			event := Event{ID: 1, Title: "Updated Event", Time: time.Now()}
			body, _ := json.Marshal(event)
			req, err := http.NewRequest("PUT", baseURL+"/events", bytes.NewReader(body))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		})

		ginkgo.It("should delete an event successfully", func() {
			req, err := http.NewRequest("DELETE", baseURL+"/events?id="+strconv.Itoa(50), nil)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			resp, err := client.Do(req)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))
		})

		ginkgo.It("should list events successfully", func() {
			req, err := http.NewRequest("GET", baseURL+"/events?startDate=2024-08-03&days=7", nil)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			resp, err := client.Do(req)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(resp.Body)

			gomega.Expect(resp.StatusCode).Should(gomega.Equal(http.StatusOK))

			var events []Event
			err = json.NewDecoder(resp.Body).Decode(&events)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())

			expectedEvents := []Event{
				{2, "Team Standup Meeting", time.Date(2024, 8, 6, 9, 0, 0, 0, time.UTC), true},
				{3, "Brunch with Clients", time.Date(2024, 8, 6, 11, 0, 0, 0, time.UTC), true},
				{4, "Product Launch Review", time.Date(2024, 8, 6, 14, 0, 0, 0, time.UTC), true},
				{5, "Workshop: Innovation Strategies", time.Date(2024, 8, 6, 16, 0, 0, 0, time.UTC), true},
				{6, "Evening Run", time.Date(2024, 8, 6, 18, 0, 0, 0, time.UTC), true},
				{7, "Dinner with Partners", time.Date(2024, 8, 6, 19, 30, 0, 0, time.UTC), true},
				{8, "Networking Event", time.Date(2024, 8, 6, 20, 30, 0, 0, time.UTC), true},
				{9, "Late Night Coding", time.Date(2024, 8, 6, 22, 0, 0, 0, time.UTC), true},
				{10, "Midnight Meditation", time.Date(2024, 8, 6, 23, 59, 0, 0, time.UTC), true},
				{11, "Morning Jogging", time.Date(2024, 8, 7, 6, 0, 0, 0, time.UTC), true},
				{12, "Team Brainstorming Session", time.Date(2024, 8, 7, 9, 0, 0, 0, time.UTC), true},
				{13, "Lunch with Investors", time.Date(2024, 8, 7, 12, 0, 0, 0, time.UTC), true},
				{14, "Design Review Meeting", time.Date(2024, 8, 7, 14, 0, 0, 0, time.UTC), true},
				{15, "Webinar: Marketing Trends", time.Date(2024, 8, 7, 16, 0, 0, 0, time.UTC), true},
				{16, "Yoga Class", time.Date(2024, 8, 7, 18, 0, 0, 0, time.UTC), true},
				{17, "Project Kickoff", time.Date(2024, 8, 7, 19, 0, 0, 0, time.UTC), true},
				{18, "Strategy Meeting", time.Date(2024, 8, 7, 20, 0, 0, 0, time.UTC), true},
				{19, "Coding Session", time.Date(2024, 8, 7, 22, 0, 0, 0, time.UTC), true},
				{20, "Late Night Reading", time.Date(2024, 8, 7, 23, 30, 0, 0, time.UTC), true},
				{21, "Morning Stretch", time.Date(2024, 8, 8, 6, 0, 0, 0, time.UTC), true},
				{22, "Scrum Meeting", time.Date(2024, 8, 8, 9, 0, 0, 0, time.UTC), true},
				{23, "Client Presentation", time.Date(2024, 8, 8, 11, 0, 0, 0, time.UTC), true},
				{24, "Product Demo", time.Date(2024, 8, 8, 13, 0, 0, 0, time.UTC), true},
				{25, "Workshop: Creative Thinking", time.Date(2024, 8, 8, 15, 0, 0, 0, time.UTC), true},
				{26, "Evening Walk", time.Date(2024, 8, 8, 18, 0, 0, 0, time.UTC), true},
				{27, "Dinner with Team", time.Date(2024, 8, 8, 19, 0, 0, 0, time.UTC), true},
				{28, "Town Hall Meeting", time.Date(2024, 8, 8, 20, 0, 0, 0, time.UTC), true},
				{29, "Evening Coding", time.Date(2024, 8, 8, 21, 0, 0, 0, time.UTC), true},
				{30, "Night Meditation", time.Date(2024, 8, 8, 22, 30, 0, 0, time.UTC), true},
				{31, "Morning Workout", time.Date(2024, 8, 9, 6, 0, 0, 0, time.UTC), true},
				{32, "Project Sync", time.Date(2024, 8, 9, 9, 0, 0, 0, time.UTC), true},
				{33, "Lunch with CEO", time.Date(2024, 8, 9, 12, 0, 0, 0, time.UTC), true},
				{34, "Design Thinking Session", time.Date(2024, 8, 9, 14, 0, 0, 0, time.UTC), true},
				{35, "Webinar: Future Tech", time.Date(2024, 8, 9, 16, 0, 0, 0, time.UTC), true},
				{36, "Evening Yoga", time.Date(2024, 8, 9, 18, 0, 0, 0, time.UTC), true},
				{37, "Dinner Meeting", time.Date(2024, 8, 9, 19, 0, 0, 0, time.UTC), true},
				{38, "Team Outing", time.Date(2024, 8, 9, 20, 0, 0, 0, time.UTC), true},
				{39, "Night Coding", time.Date(2024, 8, 9, 21, 0, 0, 0, time.UTC), true},
				{40, "Midnight Yoga", time.Date(2024, 8, 9, 23, 59, 0, 0, time.UTC), true},
			}

			gomega.Expect(events).Should(gomega.Equal(expectedEvents))
		})
	})
})

package internalhttp

import (
	"net/http"
	"time"

	"github.com/TheJubadze/OtusGolangPro/hw12_13_14_15_calendar/app/service/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		logger.Log.Infof("HTTP Request: clientIP=%s, method=%s, uri=%s, protocol=%s, statusCode=%d, contentLength=%d, userAgent=%s, duration=%dms",
			r.RemoteAddr, r.Method, r.RequestURI, r.Proto, lrw.statusCode, lrw.size, r.UserAgent(), duration.Milliseconds())
	})
}

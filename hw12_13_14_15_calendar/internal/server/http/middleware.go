package internalhttp

import (
	"log"
	"net/http"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(lrw, r)

		clientIP := r.RemoteAddr
		method := r.Method
		uri := r.RequestURI
		protocol := r.Proto
		statusCode := lrw.statusCode
		contentLength := lrw.size
		userAgent := r.UserAgent()

		log.Printf("%s %s %s %s %d %d \"%s\"\n",
			clientIP, method, uri, protocol, statusCode, contentLength, userAgent)
	})
}

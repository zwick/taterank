package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := ResponseWriteWithContext(w)

		next.ServeHTTP(rw, r)

		method := r.Method
		uri := r.URL.RequestURI()
		remoteAddr := r.RemoteAddr
		userAgent := r.UserAgent()
		proto := r.Proto
		host := r.Host
		referer := r.Referer()
		received := start.UTC()
		replied := time.Now().UTC()
		durationMS := fmt.Sprint(time.Since(start).Milliseconds())
		status := fmt.Sprint(rw.status)

		app.logger.Info(
			"request",
			"method", method,
			"uri", uri,
			"remoteAddr", remoteAddr,
			"userAgent", userAgent,
			"proto", proto,
			"host", host,
			"referer", referer,
			"status", status,
			"received", received,
			"replied", replied,
			"duration_in_ms", durationMS,
		)
	})
}

type responseWriteWithContext struct {
	writer        http.ResponseWriter
	status        int
	headerWritten bool
}

func ResponseWriteWithContext(w http.ResponseWriter) *responseWriteWithContext {
	return &responseWriteWithContext{
		writer: w,
		status: http.StatusOK,
	}
}

func (rw *responseWriteWithContext) Header() http.Header {
	return rw.writer.Header()
}

func (rw *responseWriteWithContext) WriteHeader(statusCode int) {
	rw.writer.WriteHeader(statusCode)

	if !rw.headerWritten {
		rw.status = statusCode
		rw.headerWritten = true
	}
}

func (rw *responseWriteWithContext) Write(b []byte) (int, error) {
	rw.headerWritten = true
	return rw.writer.Write(b)
}

func (rw *responseWriteWithContext) Unwrap() http.ResponseWriter {
	return rw.writer
}

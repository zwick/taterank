package main

import (
	"net/http"
)

func (app *application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		method := r.Method
		uri := r.URL.RequestURI()
		remoteAddr := r.RemoteAddr
		userAgent := r.UserAgent()
		proto := r.Proto
		host := r.Host
		referer := r.Referer()

		app.logger.Info("request", "method", method, "uri", uri, "remoteAddr", remoteAddr, "userAgent", userAgent, "proto", proto, "host", host, "referer", referer)

		next.ServeHTTP(w, r)
	})
}

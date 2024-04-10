package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"taterank.com/internal/assert"
)

var (
	buff       bytes.Buffer
	mockLogger = slog.New(slog.NewJSONHandler(&buff, nil))
)

func TestRequestLogger(t *testing.T) {
	type logOutput struct {
		Level      string `json:"level"`
		Msg        string `json:"msg"`
		Method     string `json:"method"`
		Uri        string `json:"uri"`
		RemoteAddr string `json:"remoteAddr"`
		UserAgent  string `json:"userAgent"`
		Proto      string `json:"proto"`
		Host       string `json:"host"`
		Referer    string `json:"referer"`
	}

	rr := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "/", nil)

	r.RemoteAddr = "1234"
	r.Host = "localhost:3030"
	r.Header.Set("User-Agent", "test-agent")
	r.Header.Set("Referer", "test-referer")
	r.Proto = "HTTP/1.1"

	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	app := &application{logger: mockLogger, db: nil, taters: nil}

	app.requestLogger(next).ServeHTTP(rr, r)

	expected := logOutput{
		Level:      "INFO",
		Msg:        "request",
		Method:     "GET",
		Uri:        "/",
		RemoteAddr: "1234",
		UserAgent:  "test-agent",
		Proto:      "HTTP/1.1",
		Host:       "localhost:3030",
		Referer:    "test-referer",
	}

	actual := logOutput{}
	err = json.Unmarshal(buff.Bytes(), &actual)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equals(t, expected, actual)
}

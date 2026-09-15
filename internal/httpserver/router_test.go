package httpserver

import (
	"bytes"
	"log"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func TestLCDRateLimiterBoundsRequests(t *testing.T) {
	limiter := newLCDRateLimiter(2, time.Minute)
	now := time.Unix(100, 0)
	if !limiter.Allow(now) || !limiter.Allow(now) || limiter.Allow(now) {
		t.Fatal("limiter did not enforce fixed request budget")
	}
	if !limiter.Allow(now.Add(time.Minute)) {
		t.Fatal("limiter did not reset after its window")
	}
}

func TestSanitizedLogFormatterRedactsLCDCoordinates(t *testing.T) {
	var output bytes.Buffer
	delegate := &middleware.DefaultLogFormatter{Logger: log.New(&output, "", 0), NoColor: true}
	formatter := sanitizedLogFormatter{delegate: delegate}
	request := httptest.NewRequest("GET", "/v1/lcd/space?lat=32.5&lon=-94.7", nil)
	entry := formatter.NewLogEntry(request)
	entry.Write(200, 10, nil, time.Millisecond, nil)
	logged := output.String()
	if strings.Contains(logged, "32.5") || strings.Contains(logged, "-94.7") {
		t.Fatalf("coordinates leaked into log: %s", logged)
	}
	if !strings.Contains(logged, "/v1/lcd/space?[redacted]") {
		t.Fatalf("redacted route missing from log: %s", logged)
	}
	if request.URL.RawQuery != "lat=32.5&lon=-94.7" {
		t.Fatalf("formatter mutated handler request: %s", request.URL.RawQuery)
	}
}

package spacedevs

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"ds9labs.com/space-launch-server/internal/config"
)

func TestFetchUpcomingLaunchesDoesNotRequireWritableTempDirectory(t *testing.T) {
	notDirectory := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(notDirectory, []byte("file"), 0600); err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	t.Setenv("TMPDIR", notDirectory)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"count":0,"next":null,"previous":null,"results":[]}`))
	}))
	defer server.Close()

	client := NewClient(config.SpaceDevsConfig{
		BaseURL:   server.URL,
		Timeout:   time.Second,
		UserAgent: "space-launch-server-test",
	}, nil)

	response, err := client.FetchUpcomingLaunches(context.Background(), 25)
	if err != nil {
		t.Fatalf("fetch upcoming launches: %v", err)
	}
	if response.Count != 0 || len(response.Items) != 0 {
		t.Fatalf("unexpected response: %#v", response)
	}
}

func TestFetchUpcomingLaunchesRejectsOversizedResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(make([]byte, maxRawResponseBytes+1))
	}))
	defer server.Close()
	client := NewClient(config.SpaceDevsConfig{BaseURL: server.URL, Timeout: time.Second}, nil)
	if _, err := client.FetchUpcomingLaunches(context.Background(), 1); err == nil {
		t.Fatal("oversized response accepted")
	}
}

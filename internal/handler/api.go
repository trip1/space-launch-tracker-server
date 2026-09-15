package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"ds9labs.com/space-launch-server/internal/notifications"
	"ds9labs.com/space-launch-server/internal/service"
	"ds9labs.com/space-launch-server/internal/storage"
)

type API struct {
	logger   *slog.Logger
	store    storage.Store
	upcoming service.UpcomingService
	registry *notifications.Registry
	lcdSpace lcdSpaceService
}

type lcdSpaceService interface {
	Summary(context.Context, float64, float64) service.LCDSpaceSummary
}

func NewAPI(logger *slog.Logger, store storage.Store, upcoming service.UpcomingService, registry *notifications.Registry, lcdSpace lcdSpaceService) *API {
	return &API{logger: logger, store: store, upcoming: upcoming, registry: registry, lcdSpace: lcdSpace}
}

type response struct {
	Data      any    `json:"data,omitempty"`
	Message   string `json:"message,omitempty"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func writeJSON(w http.ResponseWriter, code int, payload response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(payload); err != nil {
		w.Write([]byte(`{"status":"error","message":"failed to encode response"}`))
	}
}

func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}

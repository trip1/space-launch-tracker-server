package handler

import (
	"context"
	"net/http"
	"time"
)

func (a *API) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, response{
		Status:    "ok",
		Message:   "service is healthy",
		Timestamp: nowUTC(),
	})
}

func (a *API) Readyz(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
	defer cancel()

	if err := a.store.Ping(ctx); err != nil {
		a.logger.Warn("readiness check failed", "error", err)
		writeJSON(w, http.StatusServiceUnavailable, response{
			Status:    "error",
			Message:   "valkey unavailable",
			Timestamp: nowUTC(),
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Status:    "ok",
		Message:   "service is ready",
		Timestamp: nowUTC(),
	})
}

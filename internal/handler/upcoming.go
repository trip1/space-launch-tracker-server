package handler

import (
	"net/http"
	"strconv"

	"ds9labs.com/space-launch-server/internal/service"
)

func (a *API) UpcomingLaunches(w http.ResponseWriter, r *http.Request) {
	limit, err := readLimit(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: nowUTC(),
		})
		return
	}

	launches, err := a.upcoming.UpcomingLaunches(r.Context(), limit)
	if err != nil {
		a.logger.Error("failed to fetch upcoming launches", "error", err)
		writeJSON(w, http.StatusBadGateway, response{
			Status:    "error",
			Message:   "failed to fetch upcoming launches",
			Timestamp: nowUTC(),
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data:      launches,
		Status:    "ok",
		Timestamp: nowUTC(),
	})
}

func (a *API) UpcomingEvents(w http.ResponseWriter, r *http.Request) {
	limit, err := readLimit(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Status:    "error",
			Message:   err.Error(),
			Timestamp: nowUTC(),
		})
		return
	}

	events, err := a.upcoming.UpcomingEvents(r.Context(), limit)
	if err != nil {
		a.logger.Error("failed to fetch upcoming events", "error", err)
		writeJSON(w, http.StatusBadGateway, response{
			Status:    "error",
			Message:   "failed to fetch upcoming events",
			Timestamp: nowUTC(),
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data:      events,
		Status:    "ok",
		Timestamp: nowUTC(),
	})
}

func readLimit(r *http.Request) (int, error) {
	raw := r.URL.Query().Get("limit")
	if raw == "" {
		return 0, nil
	}

	limit, err := strconv.Atoi(raw)
	if err != nil || limit <= 0 {
		return 0, service.ErrInvalidLimit
	}

	return limit, nil
}

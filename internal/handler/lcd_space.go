package handler

import (
	"math"
	"net/http"
	"strconv"
)

func (a *API) LCDSpace(w http.ResponseWriter, r *http.Request) {
	if r.URL.RawQuery != "" {
		writeJSON(w, http.StatusBadRequest, response{Status: "error", Message: "query parameters are not accepted", Timestamp: nowUTC()})
		return
	}
	latValues := r.Header.Values("X-LCD-Latitude")
	lonValues := r.Header.Values("X-LCD-Longitude")
	if len(latValues) != 1 || len(lonValues) != 1 {
		writeJSON(w, http.StatusBadRequest, response{Status: "error", Message: "one-decimal location headers are required", Timestamp: nowUTC()})
		return
	}
	latitude, latErr := strconv.ParseFloat(latValues[0], 64)
	longitude, lonErr := strconv.ParseFloat(lonValues[0], 64)
	canonical := latErr == nil && lonErr == nil && strconv.FormatFloat(latitude, 'f', 1, 64) == latValues[0] && strconv.FormatFloat(longitude, 'f', 1, 64) == lonValues[0]
	if !canonical || math.IsNaN(latitude) || math.IsInf(latitude, 0) || math.IsNaN(longitude) || math.IsInf(longitude, 0) || latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		writeJSON(w, http.StatusBadRequest, response{Status: "error", Message: "valid lat and lon are required", Timestamp: nowUTC()})
		return
	}
	summary := a.lcdSpace.Summary(r.Context(), latitude, longitude)
	w.Header().Set("Cache-Control", "private, no-store")
	writeJSON(w, http.StatusOK, response{Status: "ok", Data: summary, Timestamp: nowUTC()})
}

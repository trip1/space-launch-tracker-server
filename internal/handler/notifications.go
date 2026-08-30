package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"ds9labs.com/space-launch-server/internal/notifications"
)

type fcmRegistrationRequest struct {
	InstallationID string `json:"installation_id"`
	Token          string `json:"token"`
	AppVersion     string `json:"app_version"`
}

type fcmRegistrationResponse struct {
	DeviceSecret string `json:"device_secret"`
}

func (a *API) RegisterFCMDevice(w http.ResponseWriter, r *http.Request) {
	if a.registry == nil {
		writeJSON(w, http.StatusServiceUnavailable, response{Status: "error", Message: "notifications unavailable", Timestamp: nowUTC()})
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	defer r.Body.Close()
	var request fcmRegistrationRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, response{Status: "error", Message: "invalid registration", Timestamp: nowUTC()})
		return
	}
	secret, err := a.registry.Register(r.Context(), notifications.Registration{InstallationID: strings.TrimSpace(request.InstallationID), Token: strings.TrimSpace(request.Token), AppVersion: strings.TrimSpace(request.AppVersion), Secret: r.Header.Get("X-Device-Secret")})
	if errors.Is(err, notifications.ErrUnauthorizedDevice) {
		writeJSON(w, http.StatusUnauthorized, response{Status: "error", Message: "invalid device credential", Timestamp: nowUTC()})
		return
	}
	if err != nil {
		a.logger.Warn("FCM registration failed", "error", err)
		writeJSON(w, http.StatusBadRequest, response{Status: "error", Message: "registration rejected", Timestamp: nowUTC()})
		return
	}
	writeJSON(w, http.StatusCreated, response{Status: "ok", Data: fcmRegistrationResponse{DeviceSecret: secret}, Timestamp: nowUTC()})
}

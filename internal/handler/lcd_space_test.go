package handler

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ds9labs.com/space-launch-server/internal/service"
)

type lcdSpaceFake struct{}

func (lcdSpaceFake) Summary(context.Context, float64, float64) service.LCDSpaceSummary {
	return service.LCDSpaceSummary{ObservedAt: 1789452000, Launch: service.LCDLaunch{Available: true, Name: "Mission", Net: 1789561800}}
}

func TestLCDSpaceEndpointReturnsCompactSummary(t *testing.T) {
	api := NewAPI(slog.New(slog.NewTextHandler(io.Discard, nil)), nil, nil, nil, lcdSpaceFake{})
	req := httptest.NewRequest(http.MethodGet, "/v1/lcd/space", nil)
	req.Header.Set("X-LCD-Latitude", "32.5")
	req.Header.Set("X-LCD-Longitude", "-94.7")
	res := httptest.NewRecorder()
	api.LCDSpace(res, req)
	if res.Code != http.StatusOK || !strings.Contains(res.Body.String(), `"name":"Mission"`) {
		t.Fatalf("status=%d body=%s", res.Code, res.Body.String())
	}
	if res.Body.Len() > 2048 {
		t.Fatalf("LCD response too large: %d", res.Body.Len())
	}
}

func TestLCDSpaceEndpointRejectsMissingOrInvalidCoordinates(t *testing.T) {
	api := NewAPI(slog.New(slog.NewTextHandler(io.Discard, nil)), nil, nil, nil, lcdSpaceFake{})
	requests := []struct{ target, lat, lon string }{
		{"/v1/lcd/space", "", ""}, {"/v1/lcd/space", "91.0", "0.0"}, {"/v1/lcd/space", "0.0", "-181.0"},
		{"/v1/lcd/space", "nan", "0.0"}, {"/v1/lcd/space", "32.55", "-94.7"}, {"/v1/lcd/space?lat=32.5&lon=-94.7", "32.5", "-94.7"},
	}
	for _, item := range requests {
		res := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, item.target, nil)
		if item.lat != "" {
			req.Header.Set("X-LCD-Latitude", item.lat)
		}
		if item.lon != "" {
			req.Header.Set("X-LCD-Longitude", item.lon)
		}
		api.LCDSpace(res, req)
		if res.Code != http.StatusBadRequest {
			t.Fatalf("%s lat=%q lon=%q status=%d", item.target, item.lat, item.lon, res.Code)
		}
	}
}

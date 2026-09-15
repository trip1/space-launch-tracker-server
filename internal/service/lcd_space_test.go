package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ds9labs.com/space-launch-server/internal/model"
)

type lcdUpcomingFake struct{}

func (lcdUpcomingFake) UpcomingLaunches(context.Context, int) (model.UpcomingLaunches, error) {
	return model.UpcomingLaunches{Items: []model.Launch{{Name: "Test Mission", Net: "2026-09-16T12:30:00Z", Status: model.LaunchStatus{Abbrev: "Go"}}}}, nil
}
func (lcdUpcomingFake) UpcomingEvents(context.Context, int) (model.UpcomingEvents, error) {
	return model.UpcomingEvents{}, nil
}

type stalledLCDUpcomingFake struct{}

func (stalledLCDUpcomingFake) UpcomingLaunches(ctx context.Context, _ int) (model.UpcomingLaunches, error) {
	<-ctx.Done()
	return model.UpcomingLaunches{}, ctx.Err()
}
func (stalledLCDUpcomingFake) UpcomingEvents(context.Context, int) (model.UpcomingEvents, error) {
	return model.UpcomingEvents{}, nil
}

func TestLCDSpaceSummaryAggregatesBoundedFeeds(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/kp", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`[{"time_tag":"2026-09-15T03:00:00","Kp":2.33},{"time_tag":"2026-09-15T06:00:00","Kp":5.00}]`))
	})
	mux.HandleFunc("/celnav", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"apiversion":"4.0.1","properties":{"data":[{"object":"Mars","almanac_data":{"hc":12.2,"zn":90.0}},{"object":"Jupiter","almanac_data":{"hc":44.8,"zn":225.0}},{"object":"VEGA","almanac_data":{"hc":80.0,"zn":1.0}}]}}`))
	})
	mux.HandleFunc("/cad", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`{"signature":{"version":"1.5"},"fields":["des","cd","dist","fullname"],"data":[["2026 AB","2026-Sep-16 03:04","0.00256955529"," (2026 AB)"]]}`))
	})
	upstream := httptest.NewServer(mux)
	defer upstream.Close()

	now := time.Date(2026, 9, 15, 6, 0, 0, 0, time.UTC)
	service := NewLCDSpaceService(lcdUpcomingFake{}, nil, LCDSpaceConfig{
		NOAAURL: upstream.URL + "/kp", USNOURL: upstream.URL + "/celnav", JPLURL: upstream.URL + "/cad", Timeout: time.Second,
	}, func() time.Time { return now })
	got := service.Summary(context.Background(), 32.5, -94.74)
	if !got.Launch.Available || got.Launch.Name != "Test Mission" || got.Launch.Net != 1789561800 {
		t.Fatalf("launch = %#v", got.Launch)
	}
	if !got.Solar.Available || got.Solar.KP != 5.0 || got.Solar.Level != "storm" {
		t.Fatalf("solar = %#v", got.Solar)
	}
	if !got.Planet.Available || got.Planet.Name != "Jupiter" || got.Planet.Altitude != 45 || got.Planet.Direction != "SW" {
		t.Fatalf("planet = %#v", got.Planet)
	}
	if !got.NEO.Available || got.NEO.Name != "2026 AB" || got.NEO.DistanceLD < 0.99 || got.NEO.DistanceLD > 1.01 {
		t.Fatalf("neo = %#v", got.NEO)
	}
	if !got.Moon.Available || got.Moon.Illumination < 0 || got.Moon.Illumination > 100 || got.Moon.Phase == "" {
		t.Fatalf("moon = %#v", got.Moon)
	}
	if !got.Mission.Available || got.Mission.Name != "Voyager 1" || got.Mission.DistanceAU < 160 {
		t.Fatalf("mission = %#v", got.Mission)
	}
}

func TestLCDSpaceSummaryDegradesIndividualFailedFeeds(t *testing.T) {
	service := NewLCDSpaceService(lcdUpcomingFake{}, nil, LCDSpaceConfig{NOAAURL: "http://127.0.0.1:1", USNOURL: "http://127.0.0.1:1", JPLURL: "http://127.0.0.1:1", Timeout: 20 * time.Millisecond}, func() time.Time { return time.Date(2026, 9, 15, 6, 0, 0, 0, time.UTC) })
	got := service.Summary(context.Background(), 0, 0)
	if !got.Launch.Available || !got.Moon.Available || !got.Mission.Available {
		t.Fatalf("local/launch feeds unavailable: %#v", got)
	}
	if got.Solar.Available || got.Planet.Available || got.NEO.Available {
		t.Fatalf("failed feeds marked available: %#v", got)
	}
}

func TestLCDSpaceSummaryBoundsStalledLaunchWithAggregateTimeout(t *testing.T) {
	service := NewLCDSpaceService(stalledLCDUpcomingFake{}, nil, LCDSpaceConfig{
		NOAAURL: "http://127.0.0.1:1", USNOURL: "http://127.0.0.1:1", JPLURL: "http://127.0.0.1:1", Timeout: 25 * time.Millisecond,
	}, func() time.Time { return time.Date(2026, 9, 15, 6, 0, 0, 0, time.UTC) })
	started := time.Now()
	got := service.Summary(context.Background(), 0, 0)
	if elapsed := time.Since(started); elapsed > 250*time.Millisecond {
		t.Fatalf("aggregate ignored timeout: %s", elapsed)
	}
	if got.Launch.Available {
		t.Fatalf("stalled launch marked available: %#v", got.Launch)
	}
}

func TestLCDSpaceSummaryJSONPreservesZeroMeasurementsAndFixedSectionKeys(t *testing.T) {
	summary := LCDSpaceSummary{Version: 1, Solar: LCDSolar{Available: true, KP: 0}, Moon: LCDMoon{Available: true, Illumination: 0}, Planet: LCDPlanet{Available: true, Altitude: 0, Azimuth: 0}}
	payload, err := json.Marshal(summary)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{`"kp":0`, `"illumination":0`, `"altitude":0`, `"azimuth":0`, `"name":""`, `"status":""`} {
		if !strings.Contains(string(payload), required) {
			t.Fatalf("fixed schema missing %s: %s", required, payload)
		}
	}
}

func TestLCDNEORejectsNonFiniteAndOutOfRangeDistances(t *testing.T) {
	for _, distance := range []string{"NaN", "+Inf", "-Inf", "0.1"} {
		t.Run(distance, func(t *testing.T) {
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Write([]byte(`{"signature":{"version":"1.5"},"fields":["des","cd","dist"],"data":[["bad","2026-Sep-16 03:04","` + distance + `"]]}`))
			}))
			defer upstream.Close()
			service := NewLCDSpaceService(lcdUpcomingFake{}, nil, LCDSpaceConfig{JPLURL: upstream.URL, Timeout: time.Second}, nil)
			if got := service.neo(context.Background()); got.Available {
				t.Fatalf("distance %q accepted: %#v", distance, got)
			}
		})
	}
}

package service

import (
	"context"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"ds9labs.com/space-launch-server/internal/storage"
)

const lunarDistanceAU = 0.00256955529

type LCDSpaceConfig struct {
	NOAAURL  string
	USNOURL  string
	JPLURL   string
	Timeout  time.Duration
	CacheTTL time.Duration
}

type LCDLaunch struct {
	Available bool   `json:"available"`
	Name      string `json:"name"`
	Net       int64  `json:"net"`
	Status    string `json:"status"`
}
type LCDMoon struct {
	Available    bool   `json:"available"`
	Phase        string `json:"phase"`
	Illumination int    `json:"illumination"`
}
type LCDSolar struct {
	Available  bool    `json:"available"`
	KP         float64 `json:"kp"`
	Level      string  `json:"level"`
	ObservedAt int64   `json:"observed_at"`
}
type LCDPlanet struct {
	Available bool   `json:"available"`
	Name      string `json:"name"`
	Altitude  int    `json:"altitude"`
	Azimuth   int    `json:"azimuth"`
	Direction string `json:"direction"`
}
type LCDNEO struct {
	Available  bool    `json:"available"`
	Name       string  `json:"name"`
	ApproachAt int64   `json:"approach_at"`
	DistanceLD float64 `json:"distance_ld"`
}
type LCDMission struct {
	Available   bool    `json:"available"`
	Name        string  `json:"name"`
	DistanceAU  float64 `json:"distance_au"`
	MissionDays int     `json:"mission_days"`
}

type LCDSpaceSummary struct {
	Version    int        `json:"v"`
	ObservedAt int64      `json:"observed_at"`
	Launch     LCDLaunch  `json:"launch"`
	Moon       LCDMoon    `json:"moon"`
	Solar      LCDSolar   `json:"solar"`
	Planet     LCDPlanet  `json:"planet"`
	NEO        LCDNEO     `json:"neo"`
	Mission    LCDMission `json:"mission"`
}

type LCDSpaceService struct {
	upcoming UpcomingService
	store    storage.Store
	config   LCDSpaceConfig
	client   *http.Client
	now      func() time.Time
	flightMu sync.Mutex
	flights  map[string]*lcdFlight
}

type lcdFlight struct {
	done   chan struct{}
	result LCDSpaceSummary
}

func NewLCDSpaceService(upcoming UpcomingService, store storage.Store, config LCDSpaceConfig, now func() time.Time) *LCDSpaceService {
	if config.NOAAURL == "" {
		config.NOAAURL = "https://services.swpc.noaa.gov/products/noaa-planetary-k-index.json"
	}
	if config.USNOURL == "" {
		config.USNOURL = "https://aa.usno.navy.mil/api/celnav"
	}
	if config.JPLURL == "" {
		config.JPLURL = "https://ssd-api.jpl.nasa.gov/cad.api"
	}
	if config.Timeout <= 0 {
		config.Timeout = 5 * time.Second
	}
	if config.CacheTTL <= 0 {
		config.CacheTTL = 5 * time.Minute
	}
	if now == nil {
		now = time.Now
	}
	return &LCDSpaceService{upcoming: upcoming, store: store, config: config, client: &http.Client{Timeout: config.Timeout}, now: now, flights: make(map[string]*lcdFlight)}
}

func (s *LCDSpaceService) Summary(ctx context.Context, latitude, longitude float64) LCDSpaceSummary {
	now := s.now().UTC()
	ctx, cancel := context.WithTimeout(ctx, s.config.Timeout)
	defer cancel()
	cacheKey := "lcd:space:" + strconv.FormatFloat(latitude, 'f', 1, 64) + ":" + strconv.FormatFloat(longitude, 'f', 1, 64)
	if s.store != nil {
		if payload, err := s.store.Get(ctx, cacheKey); err == nil {
			var cached LCDSpaceSummary
			if json.Unmarshal(payload, &cached) == nil {
				return cached
			}
		}
	}
	s.flightMu.Lock()
	if existing := s.flights[cacheKey]; existing != nil {
		s.flightMu.Unlock()
		select {
		case <-existing.done:
			return existing.result
		case <-ctx.Done():
			return LCDSpaceSummary{Version: 1, ObservedAt: now.Unix(), Moon: moonAt(now), Mission: voyagerAt(now)}
		}
	}
	flight := &lcdFlight{done: make(chan struct{})}
	s.flights[cacheKey] = flight
	s.flightMu.Unlock()
	defer func() {
		s.flightMu.Lock()
		delete(s.flights, cacheKey)
		close(flight.done)
		s.flightMu.Unlock()
	}()
	result := LCDSpaceSummary{Version: 1, ObservedAt: now.Unix(), Moon: moonAt(now), Mission: voyagerAt(now)}
	var wg sync.WaitGroup
	wg.Add(4)
	go func() { defer wg.Done(); result.Launch = s.nextLaunch(ctx) }()
	go func() { defer wg.Done(); result.Solar = s.solar(ctx) }()
	go func() { defer wg.Done(); result.Planet = s.planet(ctx, now, latitude, longitude) }()
	go func() { defer wg.Done(); result.NEO = s.neo(ctx) }()
	wg.Wait()
	if s.store != nil {
		if payload, err := json.Marshal(result); err == nil {
			_ = s.store.Set(ctx, cacheKey, payload, s.config.CacheTTL)
		}
	}
	flight.result = result
	return result
}

func (s *LCDSpaceService) nextLaunch(ctx context.Context) LCDLaunch {
	launches, err := s.upcoming.UpcomingLaunches(ctx, 1)
	if err != nil || len(launches.Items) == 0 {
		return LCDLaunch{}
	}
	launch := launches.Items[0]
	net, err := time.Parse(time.RFC3339, launch.Net)
	if err != nil {
		return LCDLaunch{}
	}
	return LCDLaunch{Available: true, Name: boundedText(launch.Name, 32), Net: net.Unix(), Status: boundedText(launch.Status.Abbrev, 12)}
}

func (s *LCDSpaceService) getJSON(ctx context.Context, endpoint string, output any) bool {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "space-launch-server/lcd")
	resp, err := s.client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return false
	}
	payload, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024+1))
	return err == nil && len(payload) <= 256*1024 && json.Unmarshal(payload, output) == nil
}

func (s *LCDSpaceService) solar(ctx context.Context) LCDSolar {
	var rows []struct {
		TimeTag string  `json:"time_tag"`
		KP      float64 `json:"Kp"`
	}
	if !s.getJSON(ctx, s.config.NOAAURL, &rows) || len(rows) == 0 {
		return LCDSolar{}
	}
	row := rows[len(rows)-1]
	observed, err := time.Parse("2006-01-02T15:04:05", row.TimeTag)
	if err != nil || row.KP < 0 || row.KP > 9 {
		return LCDSolar{}
	}
	age := s.now().UTC().Sub(observed.UTC())
	if age < -time.Hour || age > 24*time.Hour {
		return LCDSolar{}
	}
	level := "quiet"
	if row.KP >= 5 {
		level = "storm"
	} else if row.KP >= 4 {
		level = "active"
	}
	return LCDSolar{Available: true, KP: row.KP, Level: level, ObservedAt: observed.UTC().Unix()}
}

func (s *LCDSpaceService) planet(ctx context.Context, now time.Time, latitude, longitude float64) LCDPlanet {
	u, err := url.Parse(s.config.USNOURL)
	if err != nil {
		return LCDPlanet{}
	}
	query := u.Query()
	query.Set("date", now.Format("2006-01-02"))
	query.Set("time", now.Format("15:04"))
	query.Set("coords", strconv.FormatFloat(latitude, 'f', 4, 64)+","+strconv.FormatFloat(longitude, 'f', 4, 64))
	u.RawQuery = query.Encode()
	var response struct {
		APIVersion string `json:"apiversion"`
		Properties struct {
			Data []struct {
				Object  string `json:"object"`
				Almanac struct {
					Altitude float64 `json:"hc"`
					Azimuth  float64 `json:"zn"`
				} `json:"almanac_data"`
			} `json:"data"`
		} `json:"properties"`
	}
	if !s.getJSON(ctx, u.String(), &response) || !strings.HasPrefix(response.APIVersion, "4.") {
		return LCDPlanet{}
	}
	allowed := map[string]bool{"Mercury": true, "Venus": true, "Mars": true, "Jupiter": true, "Saturn": true}
	best := LCDPlanet{}
	bestAltitude := 0.0
	for _, item := range response.Properties.Data {
		if allowed[item.Object] && math.IsNaN(item.Almanac.Altitude) == false && math.IsInf(item.Almanac.Altitude, 0) == false && item.Almanac.Altitude > bestAltitude && item.Almanac.Altitude <= 90 && math.IsNaN(item.Almanac.Azimuth) == false && math.IsInf(item.Almanac.Azimuth, 0) == false && item.Almanac.Azimuth >= 0 && item.Almanac.Azimuth <= 360 {
			bestAltitude = item.Almanac.Altitude
			best = LCDPlanet{Available: true, Name: item.Object, Altitude: int(math.Round(item.Almanac.Altitude)), Azimuth: int(math.Round(item.Almanac.Azimuth)), Direction: compass(item.Almanac.Azimuth)}
		}
	}
	return best
}

func (s *LCDSpaceService) neo(ctx context.Context) LCDNEO {
	u, err := url.Parse(s.config.JPLURL)
	if err != nil {
		return LCDNEO{}
	}
	query := u.Query()
	query.Set("date-min", "now")
	query.Set("date-max", "+30")
	query.Set("dist-max", "10LD")
	query.Set("sort", "date")
	query.Set("limit", "1")
	query.Set("fullname", "true")
	u.RawQuery = query.Encode()
	var response struct {
		Signature struct {
			Version string `json:"version"`
		} `json:"signature"`
		Fields []string   `json:"fields"`
		Data   [][]string `json:"data"`
	}
	if !s.getJSON(ctx, u.String(), &response) || response.Signature.Version != "1.5" || len(response.Data) == 0 {
		return LCDNEO{}
	}
	indices := map[string]int{}
	for i, field := range response.Fields {
		indices[field] = i
	}
	row := response.Data[0]
	des, ok1 := field(row, indices, "des")
	date, ok2 := field(row, indices, "cd")
	distance, ok3 := field(row, indices, "dist")
	if !ok1 || !ok2 || !ok3 {
		return LCDNEO{}
	}
	approach, err := time.Parse("2006-Jan-02 15:04", date)
	if err != nil {
		return LCDNEO{}
	}
	au, err := strconv.ParseFloat(distance, 64)
	distanceLD := au / lunarDistanceAU
	if err != nil || math.IsNaN(au) || math.IsInf(au, 0) || au < 0 || distanceLD > 10 {
		return LCDNEO{}
	}
	return LCDNEO{Available: true, Name: boundedText(strings.TrimSpace(des), 24), ApproachAt: approach.UTC().Unix(), DistanceLD: distanceLD}
}

func moonAt(now time.Time) LCDMoon {
	reference := time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)
	cycles := now.Sub(reference).Hours() / 24 / 29.53058867
	phase := cycles - math.Floor(cycles)
	age := phase * 29.53058867
	illumination := int(math.Round((1 - math.Cos(2*math.Pi*phase)) * 50))
	name := "New"
	switch {
	case age < 1.84566:
		name = "New"
	case age < 5.53699:
		name = "Wax Crescent"
	case age < 9.22831:
		name = "First Quarter"
	case age < 12.91963:
		name = "Wax Gibbous"
	case age < 16.61096:
		name = "Full"
	case age < 20.30228:
		name = "Wan Gibbous"
	case age < 23.99361:
		name = "Last Quarter"
	case age < 27.68493:
		name = "Wan Crescent"
	}
	return LCDMoon{Available: true, Phase: name, Illumination: illumination}
}

func voyagerAt(now time.Time) LCDMission {
	launch := time.Date(1977, 9, 5, 12, 56, 0, 0, time.UTC)
	baseline := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	distance := 166.3 + now.Sub(baseline).Hours()/24/365.25*3.58
	return LCDMission{Available: true, Name: "Voyager 1", DistanceAU: distance, MissionDays: int(now.Sub(launch).Hours() / 24)}
}

func compass(azimuth float64) string {
	directions := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	index := int(math.Floor(math.Mod(azimuth+22.5, 360) / 45))
	if index < 0 {
		index += 8
	}
	return directions[index]
}
func boundedText(value string, maximum int) string {
	var b strings.Builder
	for _, r := range value {
		if b.Len() >= maximum {
			break
		}
		if r >= 0x20 && r <= 0x7e {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	return strings.TrimSpace(b.String())
}
func field(row []string, indices map[string]int, name string) (string, bool) {
	index, ok := indices[name]
	returnValue := ""
	if !ok || index < 0 || index >= len(row) {
		return returnValue, false
	}
	return row[index], true
}

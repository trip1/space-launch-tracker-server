package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"ds9labs.com/space-launch-server/internal/config"
	"ds9labs.com/space-launch-server/internal/model"
	"ds9labs.com/space-launch-server/internal/storage"
)

type fakeProvider struct {
	launches      model.UpcomingLaunches
	events        model.UpcomingEvents
	launchesCalls int
	eventsCalls   int
}

func (f *fakeProvider) FetchUpcomingLaunches(ctx context.Context, limit int) (model.UpcomingLaunches, error) {
	f.launchesCalls++
	return f.launches, nil
}

func (f *fakeProvider) FetchUpcomingEvents(ctx context.Context, limit int) (model.UpcomingEvents, error) {
	f.eventsCalls++
	return f.events, nil
}

type fakeStore struct {
	values map[string][]byte
}

func (s *fakeStore) Ping(ctx context.Context) error { return nil }

func (s *fakeStore) Get(ctx context.Context, key string) ([]byte, error) {
	if s.values == nil {
		return nil, storage.ErrCacheMiss
	}

	value, ok := s.values[key]
	if !ok {
		return nil, storage.ErrCacheMiss
	}

	return value, nil
}

func (s *fakeStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if s.values == nil {
		s.values = make(map[string][]byte)
	}

	s.values[key] = value
	return nil
}

func (s *fakeStore) Delete(ctx context.Context, key string) error               { delete(s.values, key); return nil }
func (s *fakeStore) Keys(ctx context.Context, pattern string) ([]string, error) { return nil, nil }
func (s *fakeStore) Close() error                                               { return nil }

func TestUpcomingLaunchesUsesCacheWhenPresent(t *testing.T) {
	provider := &fakeProvider{}
	store := &fakeStore{values: map[string][]byte{}}

	cached := model.UpcomingLaunches{Count: 1, Source: "cache", Items: []model.Launch{{ID: "l1", Name: "cached"}}}
	payload, err := json.Marshal(cached)
	if err != nil {
		t.Fatalf("marshal cached payload: %v", err)
	}
	store.values["spacedevs:upcoming:launches:limit=25"] = payload

	svc := &SpaceDevsUpcomingService{
		provider:             provider,
		store:                store,
		cacheTTL:             30 * time.Second,
		defaultLaunchesLimit: 25,
		defaultEventsLimit:   25,
	}

	result, err := svc.UpcomingLaunches(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.launchesCalls != 0 {
		t.Fatalf("expected provider not called on cache hit")
	}

	if result.Source != "cache" || result.Count != 1 {
		t.Fatalf("expected cached result, got %+v", result)
	}
}

func TestUpcomingLaunchesFetchesAndCachesOnMiss(t *testing.T) {
	provider := &fakeProvider{
		launches: model.UpcomingLaunches{Count: 2, Source: "thespacedevs.com", Items: []model.Launch{{ID: "l2", Name: "fetched"}}},
	}
	store := &fakeStore{}

	svc := &SpaceDevsUpcomingService{
		provider:             provider,
		store:                store,
		cacheTTL:             30 * time.Second,
		defaultLaunchesLimit: 25,
		defaultEventsLimit:   25,
	}

	result, err := svc.UpcomingLaunches(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.launchesCalls != 1 {
		t.Fatalf("expected provider called once, got %d", provider.launchesCalls)
	}

	if result.Count != 2 {
		t.Fatalf("expected fetched result count 2, got %d", result.Count)
	}

	if _, err := store.Get(context.Background(), "spacedevs:upcoming:launches:limit=25"); err != nil {
		t.Fatalf("expected cached value after fetch, got error: %v", err)
	}
}

func TestUpcomingEventsUsesDefaultLimitAndCaches(t *testing.T) {
	provider := &fakeProvider{
		events: model.UpcomingEvents{Count: 3, Source: "thespacedevs.com", Items: []model.Event{{ID: 0, Name: "event"}}},
	}
	store := &fakeStore{}

	svc := NewSpaceDevsUpcomingService(nil, store, config.SpaceDevsConfig{
		CacheTTL:      45 * time.Second,
		LaunchesLimit: 25,
		EventsLimit:   10,
	})
	svc.provider = provider

	_, err := svc.UpcomingEvents(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.eventsCalls != 1 {
		t.Fatalf("expected provider called once, got %d", provider.eventsCalls)
	}

	if _, err := store.Get(context.Background(), "spacedevs:upcoming:events:limit=10"); err != nil {
		t.Fatalf("expected cached events value, got error: %v", err)
	}
}

func TestUpcomingLaunchesIgnoreCorruptCache(t *testing.T) {
	provider := &fakeProvider{launches: model.UpcomingLaunches{Count: 1, Source: "thespacedevs.com"}}
	store := &fakeStore{values: map[string][]byte{"spacedevs:upcoming:launches:limit=25": []byte("not-json")}}

	svc := &SpaceDevsUpcomingService{
		provider:             provider,
		store:                store,
		cacheTTL:             30 * time.Second,
		defaultLaunchesLimit: 25,
		defaultEventsLimit:   25,
	}

	_, err := svc.UpcomingLaunches(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.launchesCalls != 1 {
		t.Fatalf("expected provider called when cache is corrupt")
	}
}

func TestUpcomingLaunchesFailOpenOnCacheErrors(t *testing.T) {
	provider := &fakeProvider{launches: model.UpcomingLaunches{Count: 1, Source: "thespacedevs.com"}}
	store := &errStore{}

	svc := &SpaceDevsUpcomingService{
		provider:             provider,
		store:                store,
		cacheTTL:             30 * time.Second,
		defaultLaunchesLimit: 25,
		defaultEventsLimit:   25,
	}

	_, err := svc.UpcomingLaunches(context.Background(), 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if provider.launchesCalls != 1 {
		t.Fatalf("expected provider called when cache get fails")
	}
}

type errStore struct{}

func (s *errStore) Ping(ctx context.Context) error { return nil }
func (s *errStore) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, errors.New("boom")
}
func (s *errStore) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return errors.New("boom")
}
func (s *errStore) Delete(ctx context.Context, key string) error { return errors.New("boom") }
func (s *errStore) Keys(ctx context.Context, pattern string) ([]string, error) {
	return nil, errors.New("boom")
}
func (s *errStore) Close() error { return nil }

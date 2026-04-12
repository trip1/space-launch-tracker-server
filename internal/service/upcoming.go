package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"ds9labs.com/space-launch-server/internal/config"
	"ds9labs.com/space-launch-server/internal/model"
	"ds9labs.com/space-launch-server/internal/provider/spacedevs"
	"ds9labs.com/space-launch-server/internal/storage"
)

var ErrInvalidLimit = errors.New("limit must be a positive integer")

type upcomingProvider interface {
	FetchUpcomingLaunches(ctx context.Context, limit int) (model.UpcomingLaunches, error)
	FetchUpcomingEvents(ctx context.Context, limit int) (model.UpcomingEvents, error)
}

type UpcomingService interface {
	UpcomingLaunches(ctx context.Context, limit int) (model.UpcomingLaunches, error)
	UpcomingEvents(ctx context.Context, limit int) (model.UpcomingEvents, error)
}

type SpaceDevsUpcomingService struct {
	provider             upcomingProvider
	store                storage.Store
	cacheTTL             time.Duration
	defaultLaunchesLimit int
	defaultEventsLimit   int
}

func NewSpaceDevsUpcomingService(client *spacedevs.Client, store storage.Store, cfg config.SpaceDevsConfig) *SpaceDevsUpcomingService {
	cacheTTL := cfg.CacheTTL
	if cacheTTL <= 0 {
		cacheTTL = time.Second
	}

	return &SpaceDevsUpcomingService{
		provider:             client,
		store:                store,
		cacheTTL:             cacheTTL,
		defaultLaunchesLimit: cfg.LaunchesLimit,
		defaultEventsLimit:   cfg.EventsLimit,
	}
}

func (s *SpaceDevsUpcomingService) UpcomingLaunches(ctx context.Context, limit int) (model.UpcomingLaunches, error) {
	if limit <= 0 {
		limit = s.defaultLaunchesLimit
	}
	cacheKey := fmt.Sprintf("spacedevs:upcoming:launches:limit=%d", limit)

	if payload, err := s.store.Get(ctx, cacheKey); err == nil {
		var cached model.UpcomingLaunches
		if decodeErr := json.Unmarshal(payload, &cached); decodeErr == nil {
			return cached, nil
		}
	} else if !errors.Is(err, storage.ErrCacheMiss) {
		// fail open on cache errors
	}

	result, err := s.provider.FetchUpcomingLaunches(ctx, limit)
	if err != nil {
		return model.UpcomingLaunches{}, err
	}

	if payload, marshalErr := json.Marshal(result); marshalErr == nil {
		_ = s.store.Set(ctx, cacheKey, payload, s.cacheTTL)
	}

	return result, nil
}

func (s *SpaceDevsUpcomingService) UpcomingEvents(ctx context.Context, limit int) (model.UpcomingEvents, error) {
	if limit <= 0 {
		limit = s.defaultEventsLimit
	}
	cacheKey := fmt.Sprintf("spacedevs:upcoming:events:limit=%d", limit)

	if payload, err := s.store.Get(ctx, cacheKey); err == nil {
		var cached model.UpcomingEvents
		if decodeErr := json.Unmarshal(payload, &cached); decodeErr == nil {
			return cached, nil
		}
	} else if !errors.Is(err, storage.ErrCacheMiss) {
		// fail open on cache errors
	}

	result, err := s.provider.FetchUpcomingEvents(ctx, limit)
	if err != nil {
		return model.UpcomingEvents{}, err
	}

	if payload, marshalErr := json.Marshal(result); marshalErr == nil {
		_ = s.store.Set(ctx, cacheKey, payload, s.cacheTTL)
	}

	return result, nil
}

package spacedevs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"ds9labs.com/space-launch-server/internal/config"
	"ds9labs.com/space-launch-server/internal/model"
	"ds9labs.com/space-launch-server/internal/storage"
)

type Client struct {
	baseURL    string
	userAgent  string
	httpClient *http.Client
	store      storage.Store
}

const sourceName = "thespacedevs.com"

type listResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

const rawResponseCacheTTL = 10 * time.Minute
const maxRawResponseBytes = 4 * 1024 * 1024

func NewClient(cfg config.SpaceDevsConfig, store storage.Store) *Client {
	return &Client{
		baseURL:   strings.TrimSuffix(cfg.BaseURL, "/"),
		userAgent: cfg.UserAgent,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
		store: store,
	}
}

func (c *Client) FetchUpcomingLaunches(ctx context.Context, limit int) (model.UpcomingLaunches, error) {
	var response listResponse[model.Launch]
	if err := c.getJSON(ctx, "/launches/upcoming/", limit, &response); err != nil {
		return model.UpcomingLaunches{}, err
	}

	return model.UpcomingLaunches{
		Items:  response.Results,
		Count:  response.Count,
		Source: sourceName,
	}, nil
}

func (c *Client) FetchUpcomingEvents(ctx context.Context, limit int) (model.UpcomingEvents, error) {
	var response listResponse[model.Event]
	if err := c.getJSON(ctx, "/events/upcoming/", limit, &response); err != nil {
		return model.UpcomingEvents{}, err
	}

	return model.UpcomingEvents{
		Items:  response.Results,
		Count:  response.Count,
		Source: sourceName,
	}, nil
}

func (c *Client) getJSON(ctx context.Context, path string, limit int, out any) error {
	endpoint, err := c.endpoint(path, limit)

	if err != nil {
		return err
	}

	cacheKey := c.rawBodyCacheKey(path, limit)
	if c.store != nil {
		cachedBody, cacheErr := c.store.Get(ctx, cacheKey)
		if cacheErr == nil {
			if err := json.Unmarshal(cachedBody, out); err == nil {
				return nil
			}
		} else if !errors.Is(cacheErr, storage.ErrCacheMiss) {
			// fail open on cache errors
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request spacedevs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4*1024))
		return fmt.Errorf("spacedevs status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	rawBody, err := io.ReadAll(io.LimitReader(resp.Body, maxRawResponseBytes+1))
	if err != nil {
		return fmt.Errorf("read spacedevs response: %w", err)
	}
	if len(rawBody) > maxRawResponseBytes {
		return errors.New("spacedevs response exceeds limit")
	}

	if c.store != nil {
		_ = c.store.Set(ctx, cacheKey, rawBody, rawResponseCacheTTL)
	}

	if err := json.Unmarshal(rawBody, out); err != nil {
		return fmt.Errorf("decode spacedevs response: %w", err)
	}

	return nil
}

func (c *Client) rawBodyCacheKey(path string, limit int) string {
	return fmt.Sprintf("spacedevs:raw:%s:limit=%d", strings.Trim(path, "/"), limit)
}

func (c *Client) endpoint(path string, limit int) (string, error) {
	// fmt.Printf("Parsed URL: %s\n", c.baseURL)
	// baseURL, err := url.Parse(c.baseURL)
	// if err != nil {
	// 	return "", fmt.Errorf("invalid spacedevs base url: %w", err)
	// }

	// baseURL := c.baseURL

	// ref, err := url.Parse(strings.TrimPrefix(path, "/"))
	// if err != nil {
	// 	return "", fmt.Errorf("invalid spacedevs path: %w", err)
	// }

	// joined := baseURL.ResolveReference(ref)
	// query := joined.Query()
	// if limit > 0 {
	// 	query.Set("limit", strconv.Itoa(limit))
	// }
	// joined.RawQuery = query.Encode()

	// if joined.Scheme == "" || joined.Host == "" {
	// 	return "", errors.New("invalid spacedevs endpoint")
	// }

	// return joined.String(), nil
	if limit > 0 {
		return c.baseURL + path + "?mode=detailed&format=json&limit=" + strconv.Itoa(limit), nil
	}

	return c.baseURL + path + "?mode=detailed&format=json&limit=50", nil
}

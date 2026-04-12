package model

import "time"

type Agency struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev,omitempty"`
	Country     string `json:"country,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	URL         string `json:"url,omitempty"`
}

type Location struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	TimeZone  string  `json:"time_zone,omitempty"`
}

type Pad struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Location      Location  `json:"location"`
	MapURL        string    `json:"map_url,omitempty"`
	Latitude      float64   `json:"latitude,omitempty"`
	Longitude     float64   `json:"longitude,omitempty"`
	LastUpdatedAt time.Time `json:"last_updated_at,omitempty"`
}

type ImageLicense struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Priority int    `json:"priority"`
	Link     string `json:"link,omitempty"`
}

type Image struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	ImageURL     string       `json:"image_url,omitempty"`
	ThumbnailURL string       `json:"thumbnail_url,omitempty"`
	Credit       string       `json:"credit,omitempty"`
	License      ImageLicense `json:"license"`
	SingleUse    bool         `json:"single_use"`
}

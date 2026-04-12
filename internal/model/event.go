package model

import "time"

type EventType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EventDatePrecision struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type EventInfoType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EventInfoLanguage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type EventInfo struct {
	Priority     int               `json:"priority"`
	Source       string            `json:"source"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FeatureImage *string           `json:"feature_image"`
	URL          string            `json:"url"`
	Type         EventInfoType     `json:"type"`
	Language     EventInfoLanguage `json:"language"`
}

type EventVidURL struct {
	Priority     int               `json:"priority"`
	Source       string            `json:"source"`
	Publisher    string            `json:"publisher"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FeatureImage string            `json:"feature_image"`
	URL          string            `json:"url"`
	Type         EventInfoType     `json:"type"`
	Language     EventInfoLanguage `json:"language"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	Live         bool              `json:"live"`
}

type EventUpdate struct {
	ID           int       `json:"id"`
	ProfileImage string    `json:"profile_image"`
	Comment      string    `json:"comment"`
	InfoURL      string    `json:"info_url"`
	CreatedBy    string    `json:"created_by"`
	CreatedOn    time.Time `json:"created_on"`
}

type EventAgency struct {
	ResponseMode string    `json:"response_mode"`
	ID           int       `json:"id"`
	URL          string    `json:"url"`
	Name         string    `json:"name"`
	Abbrev       string    `json:"abbrev"`
	Type         EventType `json:"type"`
}

type EventLaunch struct {
	ID               string             `json:"id"`
	URL              string             `json:"url"`
	Name             string             `json:"name"`
	ResponseMode     string             `json:"response_mode"`
	Slug             string             `json:"slug"`
	LaunchDesignator *string            `json:"launch_designator"`
	Status           LaunchStatus       `json:"status"`
	LastUpdated      time.Time          `json:"last_updated"`
	Net              time.Time          `json:"net"`
	NetPrecision     EventDatePrecision `json:"net_precision"`
	WindowEnd        time.Time          `json:"window_end"`
	WindowStart      time.Time          `json:"window_start"`
	Image            *Image             `json:"image"`
	Infographic      *string            `json:"infographic"`
}

type Event struct {
	ID            int                `json:"id"`
	URL           string             `json:"url"`
	Name          string             `json:"name"`
	InfoURLs      []EventInfo        `json:"info_urls"`
	VidURLs       []EventVidURL      `json:"vid_urls"`
	Image         *Image             `json:"image"`
	Date          time.Time          `json:"date"`
	Slug          string             `json:"slug"`
	Type          EventType          `json:"type"`
	Description   string             `json:"description"`
	WebcastLive   bool               `json:"webcast_live"`
	Location      string             `json:"location"`
	DatePrecision EventDatePrecision `json:"date_precision"`
	ResponseMode  string             `json:"response_mode"`
	Duration      string             `json:"duration"`
	Updates       []EventUpdate      `json:"updates"`
	LastUpdated   time.Time          `json:"last_updated"`
	Agencies      []EventAgency      `json:"agencies"`
	Launches      []EventLaunch      `json:"launches"`
	Expeditions   []any              `json:"expeditions"`
	Spacestations []any              `json:"spacestations"`
	Program       []LaunchProgram    `json:"program"`
	Astronauts    []any              `json:"astronauts"`
}

type UpcomingEvents struct {
	Items  []Event `json:"items"`
	Count  int     `json:"count"`
	Source string  `json:"source"`
}

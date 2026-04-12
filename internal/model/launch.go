package model

type LaunchStatus struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LaunchType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LaunchCountry struct {
	ID                      int    `json:"id"`
	Name                    string `json:"name"`
	Alpha2Code              string `json:"alpha_2_code"`
	Alpha3Code              string `json:"alpha_3_code"`
	NationalityName         string `json:"nationality_name"`
	NationalityNameComposed string `json:"nationality_name_composed"`
}

type LaunchImageLicense struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Priority int     `json:"priority"`
	Link     *string `json:"link"`
}

type LaunchImage struct {
	ID           int                `json:"id"`
	Name         string             `json:"name"`
	ImageURL     string             `json:"image_url"`
	ThumbnailURL string             `json:"thumbnail_url"`
	Credit       *string            `json:"credit"`
	License      LaunchImageLicense `json:"license"`
	SingleUse    bool               `json:"single_use"`
	Variants     []any              `json:"variants"`
}

type LaunchSocialMedia struct {
	ID   int          `json:"id"`
	Name string       `json:"name"`
	URL  *string      `json:"url"`
	Logo *LaunchImage `json:"logo"`
}

type LaunchSocialMediaLink struct {
	ID          int               `json:"id"`
	SocialMedia LaunchSocialMedia `json:"social_media"`
	URL         string            `json:"url"`
}

type LaunchAgencyList struct {
	ResponseMode string     `json:"response_mode"`
	ID           int        `json:"id"`
	URL          string     `json:"url"`
	Name         string     `json:"name"`
	Abbrev       string     `json:"abbrev"`
	Type         LaunchType `json:"type"`
}

type LaunchAgency struct {
	ResponseMode                  string                  `json:"response_mode"`
	ID                            int                     `json:"id"`
	URL                           string                  `json:"url"`
	Name                          string                  `json:"name"`
	Abbrev                        string                  `json:"abbrev"`
	Type                          LaunchType              `json:"type"`
	Featured                      bool                    `json:"featured"`
	Country                       []LaunchCountry         `json:"country"`
	Description                   string                  `json:"description"`
	Administrator                 string                  `json:"administrator"`
	FoundingYear                  int                     `json:"founding_year"`
	Launchers                     string                  `json:"launchers"`
	Spacecraft                    string                  `json:"spacecraft"`
	Parent                        any                     `json:"parent"`
	Image                         LaunchImage             `json:"image"`
	Logo                          LaunchImage             `json:"logo"`
	SocialLogo                    LaunchImage             `json:"social_logo"`
	TotalLaunchCount              int                     `json:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int                     `json:"consecutive_successful_launches"`
	SuccessfulLaunches            int                     `json:"successful_launches"`
	FailedLaunches                int                     `json:"failed_launches"`
	PendingLaunches               int                     `json:"pending_launches"`
	ConsecutiveSuccessfulLandings int                     `json:"consecutive_successful_landings"`
	SuccessfulLandings            int                     `json:"successful_landings"`
	FailedLandings                int                     `json:"failed_landings"`
	AttemptedLandings             int                     `json:"attempted_landings"`
	SuccessfulLandingsSpacecraft  int                     `json:"successful_landings_spacecraft"`
	FailedLandingsSpacecraft      int                     `json:"failed_landings_spacecraft"`
	AttemptedLandingsSpacecraft   int                     `json:"attempted_landings_spacecraft"`
	SuccessfulLandingsPayload     int                     `json:"successful_landings_payload"`
	FailedLandingsPayload         int                     `json:"failed_landings_payload"`
	AttemptedLandingsPayload      int                     `json:"attempted_landings_payload"`
	InfoURL                       string                  `json:"info_url"`
	WikiURL                       string                  `json:"wiki_url"`
	SocialMediaLinks              []LaunchSocialMediaLink `json:"social_media_links"`
}

type LaunchProgramType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LaunchMissionPatch struct {
	ID           int              `json:"id"`
	Name         string           `json:"name"`
	Priority     int              `json:"priority"`
	ImageURL     string           `json:"image_url"`
	Agency       LaunchAgencyList `json:"agency"`
	ResponseMode string           `json:"response_mode"`
}

type LaunchProgram struct {
	ResponseMode   string               `json:"response_mode"`
	ID             int                  `json:"id"`
	URL            string               `json:"url"`
	Name           string               `json:"name"`
	Image          LaunchImage          `json:"image"`
	InfoURL        *string              `json:"info_url"`
	WikiURL        *string              `json:"wiki_url"`
	Description    string               `json:"description"`
	Agencies       []LaunchAgencyList   `json:"agencies"`
	StartDate      string               `json:"start_date"`
	EndDate        *string              `json:"end_date"`
	MissionPatches []LaunchMissionPatch `json:"mission_patches"`
	Type           LaunchProgramType    `json:"type"`
}

type LaunchFamily struct {
	ResponseMode                  string         `json:"response_mode"`
	ID                            int            `json:"id"`
	Name                          string         `json:"name"`
	Manufacturer                  []LaunchAgency `json:"manufacturer"`
	Parent                        any            `json:"parent"`
	Description                   string         `json:"description"`
	Active                        bool           `json:"active"`
	MaidenFlight                  string         `json:"maiden_flight"`
	TotalLaunchCount              int            `json:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int            `json:"consecutive_successful_launches"`
	SuccessfulLaunches            int            `json:"successful_launches"`
	FailedLaunches                int            `json:"failed_launches"`
	PendingLaunches               int            `json:"pending_launches"`
	AttemptedLandings             int            `json:"attempted_landings"`
	SuccessfulLandings            int            `json:"successful_landings"`
	FailedLandings                int            `json:"failed_landings"`
	ConsecutiveSuccessfulLandings int            `json:"consecutive_successful_landings"`
}

type LaunchRocketConfiguration struct {
	ResponseMode                  string          `json:"response_mode"`
	ID                            int             `json:"id"`
	URL                           string          `json:"url"`
	Name                          string          `json:"name"`
	Families                      []LaunchFamily  `json:"families"`
	FullName                      string          `json:"full_name"`
	Variant                       string          `json:"variant"`
	Active                        bool            `json:"active"`
	IsPlaceholder                 bool            `json:"is_placeholder"`
	Manufacturer                  LaunchAgency    `json:"manufacturer"`
	Program                       []LaunchProgram `json:"program"`
	Reusable                      bool            `json:"reusable"`
	Image                         LaunchImage     `json:"image"`
	InfoURL                       string          `json:"info_url"`
	WikiURL                       string          `json:"wiki_url"`
	Description                   string          `json:"description"`
	Alias                         string          `json:"alias"`
	MinStage                      int             `json:"min_stage"`
	MaxStage                      int             `json:"max_stage"`
	Length                        float64         `json:"length"`
	Diameter                      float64         `json:"diameter"`
	MaidenFlight                  string          `json:"maiden_flight"`
	LaunchCost                    int             `json:"launch_cost"`
	LaunchMass                    float64         `json:"launch_mass"`
	LeoCapacity                   float64         `json:"leo_capacity"`
	GtoCapacity                   float64         `json:"gto_capacity"`
	GeoCapacity                   float64         `json:"geo_capacity"`
	SsoCapacity                   float64         `json:"sso_capacity"`
	ToThrust                      float64         `json:"to_thrust"`
	Apogee                        float64         `json:"apogee"`
	TotalLaunchCount              int             `json:"total_launch_count"`
	ConsecutiveSuccessfulLaunches int             `json:"consecutive_successful_launches"`
	SuccessfulLaunches            int             `json:"successful_launches"`
	FailedLaunches                int             `json:"failed_launches"`
	PendingLaunches               int             `json:"pending_launches"`
	AttemptedLandings             int             `json:"attempted_landings"`
	SuccessfulLandings            int             `json:"successful_landings"`
	FailedLandings                int             `json:"failed_landings"`
	ConsecutiveSuccessfulLandings int             `json:"consecutive_successful_landings"`
	FastestTurnaround             string          `json:"fastest_turnaround"`
}

type LaunchLauncherStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LaunchLauncher struct {
	ResponseMode       string               `json:"response_mode"`
	ID                 int                  `json:"id"`
	URL                string               `json:"url"`
	FlightProven       bool                 `json:"flight_proven"`
	SerialNumber       string               `json:"serial_number"`
	IsPlaceholder      bool                 `json:"is_placeholder"`
	Status             LaunchLauncherStatus `json:"status"`
	Image              LaunchImage          `json:"image"`
	Details            string               `json:"details"`
	SuccessfulLandings int                  `json:"successful_landings"`
	AttemptedLandings  int                  `json:"attempted_landings"`
	Flights            int                  `json:"flights"`
	LastLaunchDate     string               `json:"last_launch_date"`
	FirstLaunchDate    string               `json:"first_launch_date"`
	FastestTurnaround  string               `json:"fastest_turnaround"`
}

type LaunchCelestialBody struct {
	ResponseMode           string      `json:"response_mode"`
	ID                     int         `json:"id"`
	Name                   string      `json:"name"`
	Type                   LaunchType  `json:"type"`
	Diameter               float64     `json:"diameter"`
	Mass                   float64     `json:"mass"`
	Gravity                float64     `json:"gravity"`
	LengthOfDay            string      `json:"length_of_day"`
	Atmosphere             bool        `json:"atmosphere"`
	Image                  LaunchImage `json:"image"`
	Description            string      `json:"description"`
	WikiURL                string      `json:"wiki_url"`
	TotalAttemptedLaunches int         `json:"total_attempted_launches"`
	SuccessfulLaunches     int         `json:"successful_launches"`
	FailedLaunches         int         `json:"failed_launches"`
	TotalAttemptedLandings int         `json:"total_attempted_landings"`
	SuccessfulLandings     int         `json:"successful_landings"`
	FailedLandings         int         `json:"failed_landings"`
}

type LaunchLandingLocation struct {
	ID                 int                 `json:"id"`
	Name               string              `json:"name"`
	Active             bool                `json:"active"`
	Abbrev             string              `json:"abbrev"`
	Description        string              `json:"description"`
	Location           any                 `json:"location"`
	Longitude          *float64            `json:"longitude"`
	Latitude           *float64            `json:"latitude"`
	Image              LaunchImage         `json:"image"`
	SuccessfulLandings int                 `json:"successful_landings"`
	AttemptedLandings  int                 `json:"attempted_landings"`
	FailedLandings     int                 `json:"failed_landings"`
	CelestialBody      LaunchCelestialBody `json:"celestial_body"`
}

type LaunchLandingType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LaunchLanding struct {
	ID                int                   `json:"id"`
	URL               string                `json:"url"`
	Attempt           bool                  `json:"attempt"`
	Success           *bool                 `json:"success"`
	Description       string                `json:"description"`
	DownrangeDistance float64               `json:"downrange_distance"`
	LandingLocation   LaunchLandingLocation `json:"landing_location"`
	Type              LaunchLandingType     `json:"type"`
}

type LaunchPreviousFlight struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name"`
}

type LaunchLauncherStage struct {
	ID                   int                  `json:"id"`
	Type                 string               `json:"type"`
	Reused               bool                 `json:"reused"`
	LauncherFlightNumber int                  `json:"launcher_flight_number"`
	Launcher             LaunchLauncher       `json:"launcher"`
	PreviousFlightDate   string               `json:"previous_flight_date"`
	TurnAroundTime       string               `json:"turn_around_time"`
	Landing              LaunchLanding        `json:"landing"`
	PreviousFlight       LaunchPreviousFlight `json:"previous_flight"`
}

type LaunchRocket struct {
	ID              int                       `json:"id"`
	Configuration   LaunchRocketConfiguration `json:"configuration"`
	LauncherStage   []LaunchLauncherStage     `json:"launcher_stage"`
	SpacecraftStage []any                     `json:"spacecraft_stage"`
	Payloads        []any                     `json:"payloads"`
}

type LaunchCelestialBodyRef struct {
	ResponseMode string `json:"response_mode"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
}

type LaunchOrbit struct {
	ID            int                    `json:"id"`
	Name          string                 `json:"name"`
	Abbrev        string                 `json:"abbrev"`
	CelestialBody LaunchCelestialBodyRef `json:"celestial_body"`
}

type Mission struct {
	ID          int            `json:"id"`
	Name        string         `json:"name"`
	Type        string         `json:"type"`
	Description string         `json:"description"`
	Image       *LaunchImage   `json:"image"`
	Orbit       LaunchOrbit    `json:"orbit"`
	Agencies    []LaunchAgency `json:"agencies"`
	InfoURLs    []any          `json:"info_urls"`
	VidURLs     []any          `json:"vid_urls"`
}

type LaunchLocation struct {
	ResponseMode      string              `json:"response_mode"`
	ID                int                 `json:"id"`
	URL               string              `json:"url"`
	Name              string              `json:"name"`
	CelestialBody     LaunchCelestialBody `json:"celestial_body"`
	Active            bool                `json:"active"`
	Country           LaunchCountry       `json:"country"`
	Description       string              `json:"description"`
	Image             LaunchImage         `json:"image"`
	MapImage          string              `json:"map_image"`
	Longitude         float64             `json:"longitude"`
	Latitude          float64             `json:"latitude"`
	TimezoneName      string              `json:"timezone_name"`
	TotalLaunchCount  int                 `json:"total_launch_count"`
	TotalLandingCount int                 `json:"total_landing_count"`
}

type LaunchPad struct {
	ID                        int            `json:"id"`
	URL                       string         `json:"url"`
	Active                    bool           `json:"active"`
	Agencies                  []any          `json:"agencies"`
	Name                      string         `json:"name"`
	Image                     LaunchImage    `json:"image"`
	Description               string         `json:"description"`
	InfoURL                   *string        `json:"info_url"`
	WikiURL                   *string        `json:"wiki_url"`
	MapURL                    string         `json:"map_url"`
	Latitude                  float64        `json:"latitude"`
	Longitude                 float64        `json:"longitude"`
	Country                   LaunchCountry  `json:"country"`
	MapImage                  string         `json:"map_image"`
	TotalLaunchCount          int            `json:"total_launch_count"`
	OrbitalLaunchAttemptCount int            `json:"orbital_launch_attempt_count"`
	FastestTurnaround         string         `json:"fastest_turnaround"`
	Location                  LaunchLocation `json:"location"`
}

type LaunchUpdate struct {
	ID           int    `json:"id"`
	ProfileImage string `json:"profile_image"`
	Comment      string `json:"comment"`
	InfoURL      string `json:"info_url"`
	CreatedBy    string `json:"created_by"`
	CreatedOn    string `json:"created_on"`
}

type LaunchInfoURLType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LaunchLanguage struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type LaunchInfoURL struct {
	Priority     int               `json:"priority"`
	Source       string            `json:"source"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FeatureImage *string           `json:"feature_image"`
	URL          string            `json:"url"`
	Type         LaunchInfoURLType `json:"type"`
	Language     LaunchLanguage    `json:"language"`
}

type LaunchVidURL struct {
	Priority     int               `json:"priority"`
	Source       string            `json:"source"`
	Publisher    string            `json:"publisher"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	FeatureImage string            `json:"feature_image"`
	URL          string            `json:"url"`
	Type         LaunchInfoURLType `json:"type"`
	Language     LaunchLanguage    `json:"language"`
	StartTime    string            `json:"start_time"`
	EndTime      string            `json:"end_time"`
	Live         bool              `json:"live"`
}

type LaunchTimelineType struct {
	ID          int    `json:"id"`
	Abbrev      string `json:"abbrev"`
	Description string `json:"description"`
}

type LaunchTimelineItem struct {
	Type         LaunchTimelineType `json:"type"`
	RelativeTime string             `json:"relative_time"`
}

type Launch struct {
	ID                             string               `json:"id"`
	URL                            string               `json:"url"`
	Name                           string               `json:"name"`
	ResponseMode                   string               `json:"response_mode"`
	Slug                           string               `json:"slug"`
	LaunchDesignator               *string              `json:"launch_designator"`
	Status                         LaunchStatus         `json:"status"`
	LastUpdated                    string               `json:"last_updated"`
	Net                            string               `json:"net"`
	NetPrecision                   LaunchStatus         `json:"net_precision"`
	WindowEnd                      string               `json:"window_end"`
	WindowStart                    string               `json:"window_start"`
	Image                          LaunchImage          `json:"image"`
	Infographic                    *string              `json:"infographic"`
	Probability                    *int                 `json:"probability"`
	WeatherConcerns                *string              `json:"weather_concerns"`
	FailReason                     string               `json:"failreason"`
	Hashtag                        *string              `json:"hashtag"`
	LaunchServiceProvider          LaunchAgency         `json:"launch_service_provider"`
	Rocket                         LaunchRocket         `json:"rocket"`
	Mission                        Mission              `json:"mission"`
	Pad                            LaunchPad            `json:"pad"`
	WebcastLive                    bool                 `json:"webcast_live"`
	Program                        []LaunchProgram      `json:"program"`
	OrbitalLaunchAttemptCount      int                  `json:"orbital_launch_attempt_count"`
	LocationLaunchAttemptCount     int                  `json:"location_launch_attempt_count"`
	PadLaunchAttemptCount          int                  `json:"pad_launch_attempt_count"`
	AgencyLaunchAttemptCount       int                  `json:"agency_launch_attempt_count"`
	OrbitalLaunchAttemptCountYear  int                  `json:"orbital_launch_attempt_count_year"`
	LocationLaunchAttemptCountYear int                  `json:"location_launch_attempt_count_year"`
	PadLaunchAttemptCountYear      int                  `json:"pad_launch_attempt_count_year"`
	AgencyLaunchAttemptCountYear   int                  `json:"agency_launch_attempt_count_year"`
	FlightclubURL                  string               `json:"flightclub_url"`
	Updates                        []LaunchUpdate       `json:"updates"`
	InfoURLs                       []LaunchInfoURL      `json:"info_urls"`
	VidURLs                        []LaunchVidURL       `json:"vid_urls"`
	Timeline                       []LaunchTimelineItem `json:"timeline"`
	PadTurnaround                  string               `json:"pad_turnaround"`
	MissionPatches                 []LaunchMissionPatch `json:"mission_patches"`
}

type UpcomingLaunches struct {
	Items  []Launch `json:"items"`
	Count  int      `json:"count"`
	Source string   `json:"source"`
}

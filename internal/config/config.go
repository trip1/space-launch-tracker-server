package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultHTTPAddr           = ":8080"
	defaultReadTimeout        = 5 * time.Second
	defaultReadHeaderTimeout  = 2 * time.Second
	defaultWriteTimeout       = 10 * time.Second
	defaultIdleTimeout        = 60 * time.Second
	defaultShutdownTimeout    = 10 * time.Second
	defaultValkeyAddress      = "127.0.0.1:6379"
	defaultValkeyDialTimeout  = 2 * time.Second
	defaultValkeyReadTimeout  = 500 * time.Millisecond
	defaultValkeyWriteTimeout = 500 * time.Millisecond
	defaultValkeyPoolSize     = 50
	defaultValkeyMinIdleConns = 10
	defaultValkeyMaxRetries   = 2
	defaultValkeyPoolTimeout  = 1 * time.Second
	defaultValkeyConnMaxIdle  = 15 * time.Minute
	defaultValkeyConnMaxLife  = 2 * time.Hour
	defaultSpaceDevsBaseURL   = "https://ll.thespacedevs.com/2.3.0/"
	defaultSpaceDevsTimeout   = 15 * time.Second
	defaultSpaceDevsCacheTTL  = 10 * time.Minute
	defaultLaunchesLimit      = 50
	defaultEventsLimit        = 50
)

type Config struct {
	HTTP      HTTPConfig
	Valkey    ValkeyConfig
	SpaceDevs SpaceDevsConfig
}

type HTTPConfig struct {
	Addr              string
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
	MaxHeaderBytes    int
}

type ValkeyConfig struct {
	Address         string
	Password        string
	DB              int
	DialTimeout     time.Duration
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	PoolSize        int
	MinIdleConns    int
	MaxRetries      int
	PoolTimeout     time.Duration
	ConnMaxIdleTime time.Duration
	ConnMaxLifetime time.Duration
}

type SpaceDevsConfig struct {
	BaseURL       string
	Timeout       time.Duration
	CacheTTL      time.Duration
	UserAgent     string
	LaunchesLimit int
	EventsLimit   int
}

func Load() Config {
	return Config{
		HTTP: HTTPConfig{
			Addr:              getEnv("HTTP_ADDR", defaultHTTPAddr),
			ReadTimeout:       getEnvDuration("HTTP_READ_TIMEOUT", defaultReadTimeout),
			ReadHeaderTimeout: getEnvDuration("HTTP_READ_HEADER_TIMEOUT", defaultReadHeaderTimeout),
			WriteTimeout:      getEnvDuration("HTTP_WRITE_TIMEOUT", defaultWriteTimeout),
			IdleTimeout:       getEnvDuration("HTTP_IDLE_TIMEOUT", defaultIdleTimeout),
			ShutdownTimeout:   getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", defaultShutdownTimeout),
			MaxHeaderBytes:    getEnvInt("HTTP_MAX_HEADER_BYTES", 1<<20),
		},
		Valkey: ValkeyConfig{
			Address:         getEnv("VALKEY_ADDR", defaultValkeyAddress),
			Password:        getEnv("VALKEY_PASSWORD", ""),
			DB:              getEnvInt("VALKEY_DB", 0),
			DialTimeout:     getEnvDuration("VALKEY_DIAL_TIMEOUT", defaultValkeyDialTimeout),
			ReadTimeout:     getEnvDuration("VALKEY_READ_TIMEOUT", defaultValkeyReadTimeout),
			WriteTimeout:    getEnvDuration("VALKEY_WRITE_TIMEOUT", defaultValkeyWriteTimeout),
			PoolSize:        getEnvInt("VALKEY_POOL_SIZE", defaultValkeyPoolSize),
			MinIdleConns:    getEnvInt("VALKEY_MIN_IDLE_CONNS", defaultValkeyMinIdleConns),
			MaxRetries:      getEnvInt("VALKEY_MAX_RETRIES", defaultValkeyMaxRetries),
			PoolTimeout:     getEnvDuration("VALKEY_POOL_TIMEOUT", defaultValkeyPoolTimeout),
			ConnMaxIdleTime: getEnvDuration("VALKEY_CONN_MAX_IDLE_TIME", defaultValkeyConnMaxIdle),
			ConnMaxLifetime: getEnvDuration("VALKEY_CONN_MAX_LIFETIME", defaultValkeyConnMaxLife),
		},
		SpaceDevs: SpaceDevsConfig{
			BaseURL:       getEnv("SPACEDEVS_BASE_URL", defaultSpaceDevsBaseURL),
			Timeout:       getEnvDuration("SPACEDEVS_TIMEOUT", defaultSpaceDevsTimeout),
			CacheTTL:      getEnvDuration("SPACEDEVS_CACHE_TTL", defaultSpaceDevsCacheTTL),
			UserAgent:     getEnv("SPACEDEVS_USER_AGENT", "space-launch-server/0.1"),
			LaunchesLimit: getEnvInt("SPACEDEVS_LAUNCHES_LIMIT", defaultLaunchesLimit),
			EventsLimit:   getEnvInt("SPACEDEVS_EVENTS_LIMIT", defaultEventsLimit),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	raw := getEnv(key, "")
	if raw == "" {
		return defaultVal
	}

	parsed, err := strconv.Atoi(raw)
	if err != nil {
		return defaultVal
	}

	return parsed
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	raw := getEnv(key, "")
	if raw == "" {
		return defaultVal
	}

	parsed, err := time.ParseDuration(raw)
	if err != nil {
		return defaultVal
	}

	return parsed
}

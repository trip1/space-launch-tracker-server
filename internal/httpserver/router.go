package httpserver

import (
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ds9labs.com/space-launch-server/internal/handler"
)

type sanitizedLogFormatter struct{ delegate middleware.LogFormatter }

func (f sanitizedLogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	if r.URL.Path != "/v1/lcd/space" || r.URL.RawQuery == "" {
		return f.delegate.NewLogEntry(r)
	}
	clone := r.Clone(r.Context())
	clonedURL := *r.URL
	clone.URL = &clonedURL
	clone.RequestURI = clone.URL.EscapedPath() + "?[redacted]"
	return f.delegate.NewLogEntry(clone)
}

type lcdRateLimiter struct {
	mu          sync.Mutex
	windowStart time.Time
	window      time.Duration
	maximum     int
	count       int
}

func newLCDRateLimiter(maximum int, window time.Duration) *lcdRateLimiter {
	return &lcdRateLimiter{maximum: maximum, window: window}
}

func (l *lcdRateLimiter) Allow(now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.windowStart.IsZero() || now.Sub(l.windowStart) >= l.window {
		l.windowStart, l.count = now, 0
	}
	if l.count >= l.maximum {
		return false
	}
	l.count++
	return true
}

func (l *lcdRateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !l.Allow(time.Now()) {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewRouter(api *handler.API) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	formatter := &middleware.DefaultLogFormatter{Logger: log.New(os.Stdout, "", log.LstdFlags), NoColor: true}
	r.Use(middleware.RequestLogger(sanitizedLogFormatter{delegate: formatter}))
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Device-Secret", "X-LCD-Latitude", "X-LCD-Longitude"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         300,
	}))
	r.Use(middleware.Compress(5))
	r.Use(middleware.Timeout(15 * time.Second))

	r.Get("/healthz", api.Healthz)
	r.Get("/readyz", api.Readyz)

	r.Route("/v1", func(v1 chi.Router) {
		v1.Get("/launches/upcoming", api.UpcomingLaunches)
		v1.Get("/events/upcoming", api.UpcomingEvents)
		v1.With(newLCDRateLimiter(60, time.Minute).Middleware).Get("/lcd/space", api.LCDSpace)
		v1.Post("/devices/fcm", api.RegisterFCMDevice)
	})

	return r
}

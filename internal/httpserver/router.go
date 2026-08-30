package httpserver

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"ds9labs.com/space-launch-server/internal/handler"
)

func NewRouter(api *handler.API) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Device-Secret"},
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
		v1.Post("/devices/fcm", api.RegisterFCMDevice)
	})

	return r
}

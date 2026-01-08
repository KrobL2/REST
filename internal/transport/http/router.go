package httphandler

import (
	handler "go-rest-server/internal/transport/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter создаёт роутер с подключёнными handler’ами
func NewRouter(userHandler *handler.UserHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API маршруты
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", userHandler.HandleGetUsers)
		r.Post("/users", userHandler.HandleCreateUser)
	})

	return r
}

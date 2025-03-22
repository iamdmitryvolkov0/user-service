package routes

import (
	"github.com/go-chi/chi/v5"
	"user-srv/handlers"
	"user-srv/services"
)

func SetRoutes(userService services.UserService) *chi.Mux {
	r := chi.NewRouter()

	userHandler := handlers.NewUserHandler(userService)

	r.Post("/users", userHandler.Create)
	r.Get("/users/{id}", userHandler.GetByID)
	r.Get("/users", userHandler.GetAll)
	r.Put("/users/{id}", userHandler.Update)
	r.Delete("/users/{id}", userHandler.Delete)
	r.Post("/login", userHandler.Login)

	r.With(userHandler.AuthMiddleware).Get("/users/me", userHandler.GetCurrentUser)

	return r
}

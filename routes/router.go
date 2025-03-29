package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/swaggo/http-swagger"
	"user-srv/handlers"
	"user-srv/services"

	_ "user-srv/docs"
)

func SetRoutes(userService services.UserService) *chi.Mux {
	r := chi.NewRouter()

	userHandler := handlers.NewUserHandler(userService)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Полный URL
	))

	r.Post("/users", userHandler.Create)
	r.Get("/users/{id}", userHandler.ByID)
	r.Get("/users", userHandler.All)
	r.Put("/users/{id}", userHandler.Update)
	r.Delete("/users/{id}", userHandler.Delete)
	r.Post("/login", userHandler.Login)

	r.With(userHandler.AuthMiddleware).Get("/users/me", userHandler.CurrentUser)

	return r
}

package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"user-srv/repositories"
	"user-srv/routes"
	"user-srv/server"
	"user-srv/services"
)

func main() {
	db := services.InitDB()
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	userRepo := repositories.NewUserRepository(sqlxDB)
	userService := services.NewUserService(userRepo)

	go func() {
		router := routes.SetRoutes(userService)
		startHttpServer(router)
	}()

	server.StartGRPCServer(userService, ":50051")
}

func startHttpServer(router *chi.Mux) {
	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

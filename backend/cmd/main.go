package main

import (
	"log"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/mikheev-alexandr/pet-project/backend/internal/handlers"
	"github.com/mikheev-alexandr/pet-project/backend/internal/repository"
	"github.com/mikheev-alexandr/pet-project/backend/internal/service"
	"github.com/mikheev-alexandr/pet-project/backend/pkg/codegen"
	"github.com/mikheev-alexandr/pet-project/backend/server"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := repository.ConnectToPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("SSL_MODE"),
	})
	if err != nil {
		log.Fatalf("Failed connection to DB: %v", err)
	}
	repos := repository.NewRepository(db)

	err = codegen.Generate(repos)
	if err != nil {
		log.Fatalf("Failed generation code words: %s", err)
	}

	services := service.NewService(repos)
	validate := validator.New()
	handler := handlers.NewHandler(services, validate)

	srv := new(server.Server)
	srv.Run("8000", handler.InitRoutes())
}

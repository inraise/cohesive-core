package main

import (
	"cohesive-core/internal/db"
	"cohesive-core/internal/handler"
	"cohesive-core/internal/repository"
	"cohesive-core/internal/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Ошибка загрузки файла .env")
	}

	dbPool, err := db.NewDbPool(context.Background())
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer dbPool.Close()

	repos := repository.NewRepository(dbPool)

	authService := service.NewAuthService(repos.Auth)
	authHandler := handler.NewAuthHandler(authService)

	familyService := service.NewFamilyService(repos.Family)
	familyHandler := handler.NewFamilyHandler(familyService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// users
	router.Post("/api/v1/auth/register", authHandler.RegistrationUser)
	router.Post("/api/v1/auth/login", authHandler.LoginUser)

	// families
	router.Post("/api/v1/family", familyHandler.CreateFamily)
	router.Put("/api/v1/family", familyHandler.UpdateFamily)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Сервер успешно запущен на порту :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}

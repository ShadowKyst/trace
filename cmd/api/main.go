package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"trace/internal/delivery/http"
	"trace/internal/infrastructure/nanoid"
	"trace/internal/infrastructure/postgres"

	// "trace/internal/infrastructure/memory" // Больше не нужен по дефолту
	"trace/internal/usecase"
)

func main() {
	// 1. Config (Environment Variables)
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "trace")
	dbPass := getEnv("DB_PASSWORD", "trace_secret")
	dbName := getEnv("DB_NAME", "trace_db")

	// 2. Database Connection
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Failed to open db: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Printf("Warning: Database not ready yet (normal for first docker-compose run): %v", err)
	} else {
		log.Println("Connected to PostgreSQL")
	}

	// 3. Init Infrastructure
	repo := postgres.NewPasteRepository(db) // Используем Postgres!
	idGen := nanoid.NewGenerator()

	timeoutContext := 2 * time.Second
	uc := usecase.NewPasteUseCase(repo, idGen, timeoutContext)

	// 4. Init Echo & Server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	http.NewHandler(e, uc)

	log.Println("Starting Trace API on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

// Helper
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

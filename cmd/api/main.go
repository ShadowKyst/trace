package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/jackc/pgx/v5/stdlib"

	httpDelivery "trace/internal/delivery/http"
	"trace/internal/infrastructure/nanoid"
	"trace/internal/infrastructure/postgres"
	"trace/internal/usecase"
)

func main() {
	// 1. Config
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

	log.Println("Connecting to database...")
	for i := 0; i < 15; i++ {
		if err := db.Ping(); err == nil {
			log.Println("Successfully connected to PostgreSQL")
			break
		}
		log.Printf("Database not ready, retrying in 2s... (%d/15)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to database after retries: %v", err)
	}

	// 3. Init Infrastructure
	repo := postgres.NewPasteRepository(db)
	idGen := nanoid.NewGenerator()

	// Background Cleaner
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			if err := repo.Cleanup(); err != nil {
				log.Printf("Error cleaning up expired pastes: %v", err)
			}
		}
	}()

	timeoutContext := 2 * time.Second
	uc := usecase.NewPasteUseCase(repo, idGen, timeoutContext)

	// 4. Init Echo & Middleware
	e := echo.New()

	// [SECURITY FIX] Ручная реализация IPExtractor
	// Работает на любой версии Echo. Берет IP из заголовка Nginx.
	e.IPExtractor = func(req *http.Request) string {
		if realIP := req.Header.Get("X-Real-IP"); realIP != "" {
			return realIP
		}
		// Если заголовка нет (локальный запуск без Nginx), берем прямой IP
		return req.RemoteAddr
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))

	// Rate Limiter
	rateLimiterConfig := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{
				Rate:      20,
				Burst:     50,
				ExpiresIn: 3 * time.Minute,
			},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			return ctx.RealIP(), nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Slow down! Too many requests.",
			})
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Slow down! Too many requests.",
			})
		},
	}
	e.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))

	// 5. Init Handlers
	httpDelivery.NewHandler(e, uc)

	// 6. Start Server
	log.Println("Starting Trace API on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package main

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"trace/internal/delivery/http"
	"trace/internal/infrastructure/memory"
	"trace/internal/infrastructure/nanoid"
	"trace/internal/usecase"
)

func main() {
	// 1. Init Config (пока хардкод)
	timeoutContext := 2 * time.Second

	// 2. Init Infrastructure
	repo := memory.NewPasteRepository()
	idGen := nanoid.NewGenerator()

	// 3. Init UseCase
	uc := usecase.NewPasteUseCase(repo, idGen, timeoutContext)

	// 4. Init Echo & Middleware
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Для фронтенда на другом порту

	// 5. Init Handlers
	http.NewHandler(e, uc)

	// 6. Start Server
	log.Println("Starting Trace API on :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatal(err)
	}
}

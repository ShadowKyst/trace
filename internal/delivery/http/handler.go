package http

import (
	"net/http"
	"time"

	"trace/internal/domain"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase domain.PasteUseCase
}

func NewHandler(e *echo.Echo, uc domain.PasteUseCase) {
	h := &Handler{
		useCase: uc,
	}

	// Группа API v1
	api := e.Group("/api/v1")
	api.POST("/paste", h.CreatePaste)
	api.GET("/paste/:id", h.GetPaste)
}

// createRequest - DTO для входящего запроса
type createRequest struct {
	Content  string `json:"content"`
	Language string `json:"language"`
	TTL      int64  `json:"ttl"`      // Seconds
	Password string `json:"password"` // Optional
	Burn     bool   `json:"burn"`     // Optional
}

func (h *Handler) CreatePaste(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	ttl := time.Duration(req.TTL) * time.Second

	// Передаем новые параметры
	paste, err := h.useCase.Create(req.Content, req.Language, req.Password, ttl, req.Burn)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, paste)
}

func (h *Handler) GetPaste(c echo.Context) error {
	id := c.Param("id")
	// Пароль может прийти в Query параметре ?password=... или заголовке X-Password
	// Для простоты возьмем из Query пока
	password := c.QueryParam("password")

	paste, err := h.useCase.Get(id, password)
	if err != nil {
		if err.Error() == "paste not found or expired" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Paste not found"})
		}
		// Если нужен пароль, возвращаем 403 Forbidden но с телом ответа (чтобы фронт понял)
		if err.Error() == "password required" {
			return c.JSON(http.StatusForbidden, paste) // paste тут без контента, но с IsProtected=true
		}
		if err.Error() == "invalid password" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, paste)
}

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
	TTL      int64  `json:"ttl"` // В секундах, опционально
}

func (h *Handler) CreatePaste(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	ttl := time.Duration(req.TTL) * time.Second

	paste, err := h.useCase.Create(req.Content, req.Language, ttl)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, paste)
}

func (h *Handler) GetPaste(c echo.Context) error {
	id := c.Param("id")

	paste, err := h.useCase.Get(id)
	if err != nil {
		// В реальном проекте тут стоит проверять тип ошибки (Not Found vs Internal)
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Paste not found"})
	}

	return c.JSON(http.StatusOK, paste)
}

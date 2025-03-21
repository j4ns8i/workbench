package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"product-store/pkg/xredis"
)

type Handler struct {
	Echo   *echo.Echo
	Redis  *xredis.Client
	Logger *zerolog.Logger
}

func NewHandler(logger *zerolog.Logger, redisClient *xredis.Client) *Handler {
	e := echo.New()

	h := &Handler{
		Echo:   e,
		Redis:  &xredis.Client{UniversalClient: redisClient},
		Logger: logger,
	}
	e.GET("/healthz", h.Healthz)
	e.PUT("/product-categories", h.PutProductCategory)
	e.GET("/product-categories/:productCategoryName", h.GetProductCategory)
	e.PUT("/products", h.PutProduct)
	e.GET("/products/:productName", h.GetProduct)

	return h
}

func (h *Handler) ListenAndServe() {
	if err := h.Echo.Start(":8080"); err != nil {
		h.Logger.Fatal().Msg("Shutting down the server")
	}
}

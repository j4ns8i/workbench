package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"product-store/pkg/db"
)

type Handler struct {
	Echo   *echo.Echo
	DB     db.DB
	Logger *zerolog.Logger
}

var _ ServerInterface = (*Handler)(nil)

func NewHandler(logger *zerolog.Logger, database db.DB) *Handler {
	e := echo.New()

	h := &Handler{
		Echo:   e,
		DB:     database,
		Logger: logger,
	}

	RegisterHandlers(e, h)

	return h
}

func (h *Handler) ListenAndServe() {
	if err := h.Echo.Start(":8080"); err != nil {
		h.Logger.Fatal().Msg("Shutting down the server")
	}
}

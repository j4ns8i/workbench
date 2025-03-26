package api

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"product-store/pkg/db"
	"product-store/pkg/types"
)

func (h *Handler) Healthz(e echo.Context) error {
	if err := h.DB.CheckHealth(e.Request().Context()); err != nil {
		h.Logger.Err(err).Msg("health check on db failed")
		return e.JSON(http.StatusInternalServerError, "Redis connection error")
	}
	return e.JSON(http.StatusOK, "OK")
}

func (h *Handler) PutProductCategory(e echo.Context) error {
	var reqData types.ProductCategory
	if err := e.Bind(&reqData); err != nil {
		return e.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	ctx := e.Request().Context()
	pc, err := h.DB.PutProductCategory(ctx, reqData)
	if err != nil {
		h.Logger.Err(err).Msg("failed to put product category")
		return e.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}

	return e.JSON(http.StatusOK, pc)
}

func (h *Handler) GetProductCategory(e echo.Context, name string) error {
	logger := h.Logger.With().Str("product_category_name", name).Logger()
	pc, err := h.DB.GetProductCategory(e.Request().Context(), name)
	if err != nil {
		if errors.Is(err, db.ErrProductCategoryNotFound) {
			logger.Info().Msg("product category not found")
			return e.JSON(http.StatusNotFound, "product category not found")
		}
		logger.Err(err).Msg("error getting product category")
		return err
	}

	return e.JSON(http.StatusOK, pc)
}

func (h *Handler) PutProduct(c echo.Context) error {
	ctx := c.Request().Context()
	var reqData types.Product
	if err := c.Bind(&reqData); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	p, err := h.DB.PutProduct(ctx, reqData)
	if err != nil {
		if errors.Is(err, db.ErrProductCategoryNotFound) {
			return c.JSON(http.StatusNotFound, "product category not found")
		}
		h.Logger.Err(err).Msg("error putting product")
		return c.JSON(http.StatusInternalServerError, "unexpected error occurred")
	}

	return c.JSON(http.StatusOK, p)
}

func (h *Handler) GetProduct(e echo.Context, name string) error {
	logger := h.Logger.With().Str("product_name", name).Logger()
	product, err := h.DB.GetProduct(e.Request().Context(), name)
	if err != nil {
		if errors.Is(err, db.ErrProductNotFound) {
			logger.Info().Msg("product not found")
			return e.JSON(http.StatusNotFound, "product not found")
		}
		logger.Err(err).Msg("error getting product")
		return e.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	return e.JSON(http.StatusOK, product)
}

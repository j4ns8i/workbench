package main

import (
	"context"
	"fmt"
	"net/http"

	"product-store/pkg/api"
	"product-store/pkg/xredis"

	"github.com/cespare/xxhash/v2"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog"
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
	e.GET("/", h.GetHelloWorld)
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

func buildRedisKey(kind, name string) string {
	return kind + ":" + fmt.Sprintf("%x", xxhash.Sum64String(name))
}

func (h *Handler) GetHelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello World!")
}

func (h *Handler) Healthz(c echo.Context) error {
	// ping redis to check connection
	if err := h.Redis.Ping(c.Request().Context()).Err(); err != nil {
		h.Logger.Error().Err(err).Msg("Failed to ping Redis")
		return c.JSON(http.StatusInternalServerError, "Redis connection error")
	}
	return c.JSON(http.StatusOK, "OK")
}

func (h *Handler) PutProductCategory(c echo.Context) error {
	var category api.ProductCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	key := buildRedisKey("PRODUCTCATEGORY", category.Name)
	logger := h.Logger.With().Str("product_category_name", category.Name).Str("product_category_key", key).Logger()

	// check if the category already exists
	var (
		id  ulid.ULID
		obj xredis.ProductCategory
	)
	found, err := xredis.HGetAllScan(c.Request().Context(), h.Redis, key, &obj)
	if err != nil {
		logger.Err(err).Msg("Failed to check product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	if !found {
		// not found, create a new ID
		id = ulid.Make()
	} else {
		// Use existing ID
		u, err := api.NewULIDFromString(obj.ID)
		if err != nil {
			logger.Err(err).Msg("Failed to parse existing product category ID")
			return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
		}
		id = u
	}

	category.ID = id
	productCategoryRedis := xredis.FromAPIProductCategory(category)
	if err := h.Redis.HSet(c.Request().Context(), key, &productCategoryRedis).Err(); err != nil {
		logger.Err(err).Msg("Failed to store product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}

	return c.JSON(http.StatusOK, category)
}

func (h *Handler) GetProductCategory(c echo.Context) error {
	name := c.Param("productCategoryName")
	key := buildRedisKey("PRODUCTCATEGORY", name)
	logger := h.Logger.With().Str("product_category_name", name).Str("product_category_key", key).Logger()

	var obj xredis.ProductCategory
	found, err := xredis.HGetAllScan(c.Request().Context(), h.Redis, key, &obj)
	if err != nil {
		logger.Err(err).Msg("Failed to retrieve product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	if !found {
		return c.JSON(http.StatusNotFound, "Product category not found")
	}

	category, err := xredis.ToAPIProductCategory(obj)
	if err != nil {
		logger.Err(err).Msg("Failed to convert product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}

	return c.JSON(http.StatusOK, category)
}

func (h *Handler) PutProduct(c echo.Context) error {
	ctx := c.Request().Context()
	var product api.Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	t := xredis.NewTransaction(h.Redis)
	t.Prepare(xredis.WithProductCategoryExists(product.Category))
	err := t.Exec(ctx, func(ctx context.Context, tx *xredis.Tx) error {
		// Check if the product already exists to preserve its ID
		var exists = true
		obj, err := tx.GetProduct(ctx, product.Name)
		if err != nil {
			if err == xredis.ErrorNotFound {
				exists = false
			} else {
				return err
			}
		}
		if !exists {
			product.ID = ulid.Make()
		} else {
			product.ID = obj.ID
		}

		return tx.PutProduct(ctx, product)
	})

	if err != nil {
		if err == xredis.ErrorProductCategoryNotFound {
			return c.JSON(http.StatusNotFound, "Product category not found")
		}
		h.Logger.Err(err).Msg("Failed to execute transaction")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}

	return c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProduct(c echo.Context) error {
	name := c.Param("productName")
	key := buildRedisKey("PRODUCT", name)
	logger := h.Logger.With().Str("product_name", name).Str("product_key", key).Logger()

	var obj xredis.Product
	found, err := xredis.HGetAllScan(c.Request().Context(), h.Redis, key, &obj)
	if err != nil {
		logger.Err(err).Msg("failed to retrieve product")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	if !found {
		logger.Info().Msg("product not found")
		return c.JSON(http.StatusNotFound, "Product not found")
	}

	product, err := xredis.ToAPIProduct(obj)
	if err != nil {
		logger.Err(err).Msg("Failed to convert product")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	return c.JSON(http.StatusOK, product)
}

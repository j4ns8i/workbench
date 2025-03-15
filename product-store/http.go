package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cespare/xxhash/v2"
	"github.com/labstack/echo/v4"
	"github.com/oklog/ulid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type Handler struct {
	Echo   *echo.Echo
	Redis  *redis.Client
	Logger zerolog.Logger
}

func NewHandler(redisClient *redis.Client) *Handler {
	e := echo.New()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	h := &Handler{Echo: e, Redis: redisClient, Logger: logger}
	e.GET("/", h.GetHelloWorld)
	e.GET("/healthz", h.Healthz)
	e.PUT("/product-categories", h.PutProductCategory)
	e.GET("/product-categories/:productCategoryName", h.GetProductCategory)
	e.PUT("/products", h.PutProduct)
	e.GET("/product/:productName", h.GetProduct)

	return h
}

func (h *Handler) ListenAndServe() {
	if err := h.Echo.Start(":8080"); err != nil {
		h.Echo.Logger.Fatal("Shutting down the server")
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
	var category ProductCategory
	if err := c.Bind(&category); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request payload")
	}

	key := fmt.Sprintf("PRODUCTCATEGORY:%s", category.Name)

	// Check if the category already exists
	var val ProductCategoryRedis
	err := h.Redis.HGetAll(c.Request().Context(), key).Scan(&val)
	if err == redis.Nil {
		// Generate new ULID if category doesn't exist
		category.ID = ulid.Make()
	} else if err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to check product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	} else {
		u, err := NewULIDFromString(val.ID)
		if err != nil {
			h.Logger.Error().Err(err).Str("key", key).Msg("Failed to parse existing product category ID")
			return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
		}
		category.ID = u
	}

	category.ID = ulid.Make()
	productCategoryRedis := RedisFromProductCategory(category)
	if err := h.Redis.HSet(c.Request().Context(), key, &productCategoryRedis).Err(); err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to store product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}

	return c.JSON(http.StatusOK, category)
}

func (h *Handler) GetProductCategory(c echo.Context) error {
	name := c.Param("productCategoryName")
	key := buildRedisKey("PRODUCTCATEGORY", name)
	val, err := h.Redis.Get(c.Request().Context(), key).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, "Product category not found")
	} else if err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to retrieve product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	var category ProductCategory
	if err := json.Unmarshal([]byte(val), &category); err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to unmarshal product category")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	return c.JSON(http.StatusOK, category)
}

func (h *Handler) PutProduct(c echo.Context) error {
	var product Product
	if err := c.Bind(&product); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	key := buildRedisKey("PRODUCT", product.Name)

	err := h.Redis.Watch(c.Request().Context(), func(tx *redis.Tx) error {
		val, err := tx.Get(c.Request().Context(), key).Result()
		if err != nil && err != redis.Nil {
			return err
		}

		if err == redis.Nil {
			product.ID = NewULID()
		} else {
			var existingProduct Product
			if err := json.Unmarshal([]byte(val), &existingProduct); err != nil {
				return err
			}
			product.ID = existingProduct.ID
		}

		_, err = tx.TxPipelined(c.Request().Context(), func(pipe redis.Pipeliner) error {
			return pipe.Set(c.Request().Context(), key, product, 0).Err()
		})
		return err
	}, key)

	if err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to store product")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	return c.JSON(http.StatusOK, product)
}

func (h *Handler) GetProduct(c echo.Context) error {
	name := c.Param("productName")
	key := buildRedisKey("PRODUCT", name)
	val, err := h.Redis.Get(c.Request().Context(), key).Result()
	if err == redis.Nil {
		return c.JSON(http.StatusNotFound, "Product not found")
	} else if err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to retrieve product")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	var product Product
	if err := json.Unmarshal([]byte(val), &product); err != nil {
		h.Logger.Error().Err(err).Str("key", key).Msg("Failed to unmarshal product")
		return c.JSON(http.StatusInternalServerError, "Unexpected error occurred")
	}
	return c.JSON(http.StatusOK, product)
}

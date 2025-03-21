package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"product-store/pkg/api"
	"product-store/pkg/types"
	"product-store/pkg/xredis"
)

const epsilon = 0.0000001

// TestAPISuite is a test suite for API integration tests
type TestAPISuite struct {
	suite.Suite
	handler *api.Handler
	client  *http.Client
	server  *httptest.Server
	redis   *redis.Client
}

var (
	PS_TEST_REDIS_HOST     = os.Getenv("PS_TEST_REDIS_HOST")
	PS_TEST_REDIS_PORT     = os.Getenv("PS_TEST_REDIS_PORT")
	PS_TEST_REDIS_PASSWORD = os.Getenv("PS_TEST_REDIS_PASSWORD")
)

// SetupSuite sets up the test suite
func (s *TestAPISuite) SetupSuite() {
	if !EnableIntegrationTests {
		s.T().Skip("Integration tests disabled; enable by re-running task with INTEGRATION=1")
	}

	// Setup Redis client
	s.redis = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", PS_TEST_REDIS_HOST, PS_TEST_REDIS_PORT),
		Password:    PS_TEST_REDIS_PASSWORD,
		DB:          1,
		Protocol:    3,
		MaxRetries:  10,
		DialTimeout: 1 * time.Second,
	})

	// Ping Redis to ensure it's available
	_, err := s.redis.Ping(context.Background()).Result()
	s.Require().NoError(err, "Redis must be available for integration tests")

	// Create logger
	logger := log.Logger

	// Create handler
	redisClient := &xredis.Client{UniversalClient: s.redis}
	s.handler = api.NewHandler(&logger, redisClient)

	// Create test server
	s.server = httptest.NewServer(s.handler.Echo)
	s.client = &http.Client{
		Timeout: 5 * time.Second,
	}
}

// TearDownSuite closes resources after all tests
func (s *TestAPISuite) TearDownSuite() {
	// Close test server
	s.server.Close()

	// Close Redis client
	err := s.redis.Close()
	s.Require().NoError(err, "Failed to close Redis client")
}

// TearDownTest cleans up after each test
func (s *TestAPISuite) TearDownTest() {
	// Clear test data from Redis
	keys, err := s.redis.Keys(context.Background(), "*").Result()
	s.Require().NoError(err, "Failed to get Redis keys")

	if len(keys) > 0 {
		_, err = s.redis.Del(context.Background(), keys...).Result()
		s.Require().NoError(err, "Failed to clear Redis test data")
	}
}

// TestHealthz tests the healthz endpoint
func (s *TestAPISuite) TestHealthz() {
	resp, err := s.client.Get(s.server.URL + "/healthz")
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	// Verify response
	s.Equal(http.StatusOK, resp.StatusCode)

	var responseBody string
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	s.Require().NoError(err, "Failed to decode response")
	s.Equal("OK", responseBody)
}

// TestProductCategoryCRUD tests creating and retrieving a product category
func (s *TestAPISuite) TestProductCategoryCRUD() {
	// Test data
	categoryName := "Electronics"
	category := types.ProductCategory{
		ProductCategoryData: types.ProductCategoryData{
			Name: categoryName,
		},
	}

	// Create product category
	jsonData, err := json.Marshal(category)
	s.Require().NoError(err, "Failed to marshal product category")

	req, err := http.NewRequest(http.MethodPut, s.server.URL+"/product-categories", strings.NewReader(string(jsonData)))
	s.Require().NoError(err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var createdCategory types.ProductCategory
	err = json.NewDecoder(resp.Body).Decode(&createdCategory)
	s.Require().NoError(err, "Failed to decode response")
	s.Equal(categoryName, createdCategory.Name)
	s.NotEmpty(createdCategory.ID, "ID should not be empty")

	// Get product category
	resp, err = s.client.Get(fmt.Sprintf("%s/product-categories/%s", s.server.URL, categoryName))
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var retrievedCategory types.ProductCategory
	err = json.NewDecoder(resp.Body).Decode(&retrievedCategory)
	s.Require().NoError(err, "Failed to decode response")
	s.Equal(createdCategory.ID, retrievedCategory.ID)
	s.Equal(categoryName, retrievedCategory.Name)

	// Test retrieving non-existent category
	resp, err = s.client.Get(fmt.Sprintf("%s/product-categories/nonexistent", s.server.URL))
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestProductCRUD tests creating and retrieving a product
func (s *TestAPISuite) TestProductCRUD() {
	// First create a product category
	categoryName := "Electronics"
	category := types.ProductCategory{
		ProductCategoryData: types.ProductCategoryData{
			Name: categoryName,
		},
	}

	jsonData, err := json.Marshal(category)
	s.Require().NoError(err, "Failed to marshal product category")

	req, err := http.NewRequest(http.MethodPut, s.server.URL+"/product-categories", strings.NewReader(string(jsonData)))
	s.Require().NoError(err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()
	s.Equal(http.StatusOK, resp.StatusCode)

	// Create product
	productName := "Laptop"
	product := types.Product{
		ProductData: types.ProductData{
			Name:     productName,
			Category: categoryName,
			Price:    999.99,
		},
	}

	jsonData, err = json.Marshal(product)
	s.Require().NoError(err, "Failed to marshal product")

	req, err = http.NewRequest(http.MethodPut, s.server.URL+"/products", strings.NewReader(string(jsonData)))
	s.Require().NoError(err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var createdProduct types.Product
	err = json.NewDecoder(resp.Body).Decode(&createdProduct)
	s.Require().NoError(err, "Failed to decode response")
	s.Equal(productName, createdProduct.Name)
	s.Equal(categoryName, createdProduct.Category)
	s.InEpsilon(999.99, createdProduct.Price, epsilon)
	s.NotEmpty(createdProduct.ID, "ID should not be empty")

	// Get product
	resp, err = s.client.Get(fmt.Sprintf("%s/products/%s", s.server.URL, productName))
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var retrievedProduct types.Product
	err = json.NewDecoder(resp.Body).Decode(&retrievedProduct)
	s.Require().NoError(err, "Failed to decode response")
	s.Equal(createdProduct.ID, retrievedProduct.ID)
	s.Equal(productName, retrievedProduct.Name)
	s.Equal(categoryName, retrievedProduct.Category)
	s.InEpsilon(999.99, retrievedProduct.Price, epsilon)

	// Test retrieving non-existent product
	resp, err = s.client.Get(fmt.Sprintf("%s/products/nonexistent", s.server.URL))
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestProductWithNonExistentCategory tests creating a product with non-existent category
func (s *TestAPISuite) TestProductWithNonExistentCategory() {
	product := types.Product{
		ProductData: types.ProductData{
			Name:     "Invalid Product",
			Category: "NonExistentCategory",
			Price:    10.99,
		},
	}

	jsonData, err := json.Marshal(product)
	s.Require().NoError(err, "Failed to marshal product")

	req, err := http.NewRequest(http.MethodPut, s.server.URL+"/products", strings.NewReader(string(jsonData)))
	s.Require().NoError(err, "Failed to create request")
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err, "Failed to send request")
	defer resp.Body.Close()

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

// TestIdempotentCategoryCreation tests that creating the same category twice preserves the ID
func (s *TestAPISuite) TestIdempotentCategoryCreation() {
	categoryName := "Furniture"
	category := types.ProductCategory{
		ProductCategoryData: types.ProductCategoryData{
			Name: categoryName,
		},
	}

	// First creation
	jsonData, err := json.Marshal(category)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPut, s.server.URL+"/product-categories", strings.NewReader(string(jsonData)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var firstCategory types.ProductCategory
	err = json.NewDecoder(resp.Body).Decode(&firstCategory)
	s.Require().NoError(err)

	// Second creation
	req, err = http.NewRequest(http.MethodPut, s.server.URL+"/product-categories", strings.NewReader(string(jsonData)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	var secondCategory types.ProductCategory
	err = json.NewDecoder(resp.Body).Decode(&secondCategory)
	s.Require().NoError(err)

	s.Equal(firstCategory.ID, secondCategory.ID, "ID should be preserved on second creation")
}

// TestIdempotentProductCreation tests that creating the same product twice preserves the ID
func (s *TestAPISuite) TestIdempotentProductCreation() {
	// Create category first
	categoryName := "Books"
	category := types.ProductCategory{
		ProductCategoryData: types.ProductCategoryData{
			Name: categoryName,
		},
	}

	jsonData, err := json.Marshal(category)
	s.Require().NoError(err)

	req, err := http.NewRequest(http.MethodPut, s.server.URL+"/product-categories", strings.NewReader(string(jsonData)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	// Create product
	productName := "Programming Book"
	product := types.Product{
		ProductData: types.ProductData{
			Name:     productName,
			Category: categoryName,
			Price:    29.99,
		},
	}

	jsonData, err = json.Marshal(product)
	s.Require().NoError(err)

	req, err = http.NewRequest(http.MethodPut, s.server.URL+"/products", strings.NewReader(string(jsonData)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.Require().Equal(http.StatusOK, resp.StatusCode)
	s.Require().NoError(err)
	defer resp.Body.Close()

	var firstProduct types.Product
	err = json.NewDecoder(resp.Body).Decode(&firstProduct)
	s.Require().NoError(err)

	// Create product again
	req, err = http.NewRequest(http.MethodPut, s.server.URL+"/products", strings.NewReader(string(jsonData)))
	s.Require().NoError(err)
	req.Header.Set("Content-Type", "application/json")

	resp, err = s.client.Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	var secondProduct types.Product
	err = json.NewDecoder(resp.Body).Decode(&secondProduct)
	s.Require().NoError(err)

	s.Equal(firstProduct.ID, secondProduct.ID, "ID should be preserved on second creation")
}

// TestAPI runs the test suite
func TestAPI(t *testing.T) {
	suite.Run(t, new(TestAPISuite))
}

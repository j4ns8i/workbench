package integration

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"

	"product-store/pkg/api"
	"product-store/pkg/client"
	"product-store/pkg/types"
	"product-store/pkg/xredis"
)

const epsilon = 0.0000001

// TestAPISuite is a test suite for API integration tests
type TestAPISuite struct {
	suite.Suite
	server      *httptest.Server
	redisClient *redis.Client
	psClient    client.ClientWithResponsesInterface

	// test-specific variables
	id string
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
	s.redisClient = redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%s", PS_TEST_REDIS_HOST, PS_TEST_REDIS_PORT),
		Password:    PS_TEST_REDIS_PASSWORD,
		DB:          1, // Use a different DB for testing
		Protocol:    3,
		MaxRetries:  10,
		DialTimeout: 1 * time.Second,
	})

	// Ping Redis to ensure it's available
	_, err := s.redisClient.Ping(context.Background()).Result()
	s.Require().NoError(err, "Redis must be available for integration tests")

	// Create logger
	logger := log.Logger

	// Create handler
	rdb := xredis.NewDB(s.redisClient, &logger)
	handler := api.NewHandler(&logger, rdb)

	// Create test server
	s.server = httptest.NewServer(handler.Echo)

	s.psClient, err = client.NewClientWithResponses(s.server.URL)
	s.Require().NoError(err)
}

func (s *TestAPISuite) SetupTest() {
	s.id = types.NewULID().String()
}

// withID returns token with a ULID suffix to ensure uniqueness across tests.
func (s *TestAPISuite) withID(token string) string {
	return fmt.Sprintf("%s-%s", token, s.id)
}

// TearDownSuite closes resources after all tests
func (s *TestAPISuite) TearDownSuite() {
	// Close test server
	s.server.Close()

	err := s.redisClient.FlushDB(s.T().Context()).Err()
	s.NoError(err)

	// Close Redis client
	err = s.redisClient.Close()
	s.Require().NoError(err, "Failed to close Redis client")
}

// TestHealthz tests the healthz endpoint
func (s *TestAPISuite) TestHealthz() {
	ctx := s.T().Context()
	resp, err := s.psClient.HealthzWithResponse(ctx)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, resp.StatusCode())
	s.Require().NotNil(resp.JSON200)
	s.Equal("OK", *resp.JSON200)
}

// TestProductCategoryCRUD tests creating and retrieving a product category
func (s *TestAPISuite) TestProductCategoryCRUD() {
	ctx := s.T().Context()
	category := types.ProductCategory{
		Name: s.withID("Electronics"),
	}
	putPCResponse, err := s.psClient.PutProductCategoryWithResponse(ctx, category)

	// Test that creating the product category works.
	s.Require().NoError(err)
	s.Equal(http.StatusOK, putPCResponse.StatusCode())
	s.Require().NotNil(putPCResponse.JSON200)
	createdPC := *putPCResponse.JSON200
	s.Equal(category.Name, createdPC.Name)
	s.NotEmpty(createdPC.ID)

	// Test that getting the product category returns the same object we got
	// when creating the category.
	getPCResponse, err := s.psClient.GetProductCategoryWithResponse(ctx, category.Name)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, getPCResponse.StatusCode())
	s.Require().NotNil(getPCResponse.JSON200)

	gotProductCategory := *getPCResponse.JSON200
	s.Equal(createdPC.ID, gotProductCategory.ID)
	s.Equal(createdPC.Name, gotProductCategory.Name)

	// Test retrieving non-existent category.
	notFoundPCResponse, err := s.psClient.GetProductCategoryWithResponse(ctx, "not-found")
	s.Require().NoError(err)
	s.Equal(http.StatusNotFound, notFoundPCResponse.StatusCode())
}

// TestProductCRUD tests creating and retrieving a product
func (s *TestAPISuite) TestProductCRUD() {
	ctx := s.T().Context()

	// First create a product category
	category := types.ProductCategory{
		Name: s.withID("Electronics"),
	}

	putPCResponse, err := s.psClient.PutProductCategoryWithResponse(ctx, category)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, putPCResponse.StatusCode())

	// Create product
	product := types.Product{
		Name:     s.withID("Laptop"),
		Category: category.Name,
		Price:    999.99,
	}

	putProductResponse, err := s.psClient.PutProductWithResponse(ctx, product)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, putProductResponse.StatusCode())
	s.Require().NotNil(putProductResponse.JSON200)

	createdProduct := *putProductResponse.JSON200
	s.Equal(product.Name, createdProduct.Name)
	s.Equal(category.Name, createdProduct.Category)
	s.InEpsilon(999.99, createdProduct.Price, epsilon)
	s.NotEmpty(createdProduct.ID)

	// Get product
	getProductResponse, err := s.psClient.GetProductWithResponse(ctx, product.Name)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, getProductResponse.StatusCode())
	s.Require().NotNil(getProductResponse.JSON200)

	retrievedProduct := *getProductResponse.JSON200
	s.Equal(createdProduct.ID, retrievedProduct.ID)
	s.Equal(product.Name, retrievedProduct.Name)
	s.Equal(category.Name, retrievedProduct.Category)
	s.InEpsilon(999.99, retrievedProduct.Price, epsilon)

	// Test retrieving non-existent product
	nonExistentResponse, err := s.psClient.GetProductWithResponse(ctx, "nonexistent")
	s.Require().NoError(err)
	s.Equal(http.StatusNotFound, nonExistentResponse.StatusCode())
}

// TestProductWithNonExistentCategory tests creating a product with a non-existent category
func (s *TestAPISuite) TestProductWithNonExistentCategory() {
	ctx := s.T().Context()

	product := types.Product{
		Name:     s.withID("Invalid Product"),
		Category: s.withID("NonExistentCategory"),
		Price:    123.45,
	}

	response, err := s.psClient.PutProductWithResponse(ctx, product)
	s.Require().NoError(err)
	s.Equal(http.StatusNotFound, response.StatusCode())
	s.Require().NotNil(response.JSON404)
	s.Contains(*response.JSON404, "product category not found")
}

// TestIdempotentCategoryCreation tests that creating the same category twice works
func (s *TestAPISuite) TestIdempotentCategoryCreation() {
	ctx := s.T().Context()

	category := types.ProductCategory{
		Name: s.withID("IdempotentCategory"),
	}

	// Create category first time
	response1, err := s.psClient.PutProductCategoryWithResponse(ctx, category)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response1.StatusCode())
	s.Require().NotNil(response1.JSON200)
	s.Equal(category.Name, response1.JSON200.Name)

	// Create the same category again
	response2, err := s.psClient.PutProductCategoryWithResponse(ctx, category)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response2.StatusCode())
	s.Require().NotNil(response2.JSON200)
	s.Equal(category.Name, response2.JSON200.Name)
}

// TestIdempotentProductCreation tests that creating the same product twice works
func (s *TestAPISuite) TestIdempotentProductCreation() {
	ctx := s.T().Context()

	// First create a product category
	category := types.ProductCategory{
		Name: s.withID("IdempotentProductCategory"),
	}

	catResponse, err := s.psClient.PutProductCategoryWithResponse(ctx, category)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, catResponse.StatusCode())

	// Create a product
	product := types.Product{
		Name:     s.withID("IdempotentProduct"),
		Category: category.Name,
		Price:    456.78,
	}

	// Create product first time
	response1, err := s.psClient.PutProductWithResponse(ctx, product)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response1.StatusCode())
	s.Require().NotNil(response1.JSON200)
	firstID := response1.JSON200.ID
	s.NotEmpty(firstID)

	// Create the same product again
	response2, err := s.psClient.PutProductWithResponse(ctx, product)
	s.Require().NoError(err)
	s.Equal(http.StatusOK, response2.StatusCode())
	s.Require().NotNil(response2.JSON200)
	s.Equal(firstID, response2.JSON200.ID)
	s.Equal(product.Name, response2.JSON200.Name)
	s.Equal(category.Name, response2.JSON200.Category)
	s.InEpsilon(456.78, response2.JSON200.Price, epsilon)
}

// TestAPI runs the test suite
func TestAPI(t *testing.T) {
	suite.Run(t, new(TestAPISuite))
}

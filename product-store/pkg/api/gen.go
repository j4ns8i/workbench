// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	externalRef0 "product-store/pkg/types"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// PutProductCategoryJSONRequestBody defines body for PutProductCategory for application/json ContentType.
type PutProductCategoryJSONRequestBody = externalRef0.ProductCategory

// PutProductJSONRequestBody defines body for PutProduct for application/json ContentType.
type PutProductJSONRequestBody = externalRef0.Product

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get product events
	// (GET /events/products)
	GetEventsProducts(ctx echo.Context) error
	// Health check endpoint
	// (GET /healthz)
	Healthz(ctx echo.Context) error
	// Create or update a product category
	// (PUT /product-categories)
	PutProductCategory(ctx echo.Context) error
	// Get a product category by name
	// (GET /product-categories/{productCategoryName})
	GetProductCategory(ctx echo.Context, productCategoryName string) error
	// Create or update a product
	// (PUT /products)
	PutProduct(ctx echo.Context) error
	// Get a product by name
	// (GET /products/{productName})
	GetProduct(ctx echo.Context, productName string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetEventsProducts converts echo context to params.
func (w *ServerInterfaceWrapper) GetEventsProducts(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetEventsProducts(ctx)
	return err
}

// Healthz converts echo context to params.
func (w *ServerInterfaceWrapper) Healthz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Healthz(ctx)
	return err
}

// PutProductCategory converts echo context to params.
func (w *ServerInterfaceWrapper) PutProductCategory(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutProductCategory(ctx)
	return err
}

// GetProductCategory converts echo context to params.
func (w *ServerInterfaceWrapper) GetProductCategory(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "productCategoryName" -------------
	var productCategoryName string

	err = runtime.BindStyledParameterWithOptions("simple", "productCategoryName", ctx.Param("productCategoryName"), &productCategoryName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter productCategoryName: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProductCategory(ctx, productCategoryName)
	return err
}

// PutProduct converts echo context to params.
func (w *ServerInterfaceWrapper) PutProduct(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutProduct(ctx)
	return err
}

// GetProduct converts echo context to params.
func (w *ServerInterfaceWrapper) GetProduct(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "productName" -------------
	var productName string

	err = runtime.BindStyledParameterWithOptions("simple", "productName", ctx.Param("productName"), &productName, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter productName: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetProduct(ctx, productName)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/events/products", wrapper.GetEventsProducts)
	router.GET(baseURL+"/healthz", wrapper.Healthz)
	router.PUT(baseURL+"/product-categories", wrapper.PutProductCategory)
	router.GET(baseURL+"/product-categories/:productCategoryName", wrapper.GetProductCategory)
	router.PUT(baseURL+"/products", wrapper.PutProduct)
	router.GET(baseURL+"/products/:productName", wrapper.GetProduct)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWS4/bNhD+K8S0R8VymxQIdEvTIjVatMamPQWLgkuNbKYSyQyH7rqG/3tBSn5Ku/Eq",
	"r80pRjQcznwPfrsBZRtnDRr2UGzAqyU2Mv3ktUP/95xsGRTH/3BkHRJrTJ9fSsaFpXX8XaJXpB1ra6CA",
	"P5cojGxQ2ErwEoVrWwi1O5Gl3lCAZ9JmAdsMZj8N9wlGvwsodImGdaWRRGXpuCtkUFlqJEMBodblUPPf",
	"ZYMXjznUYE5a3dHBxU/9FvuZqtrKo54mNDdIsN1mQPguaMISijfthNkB092V1/uD9uYtKo7DnPByzMIp",
	"P6MRPebp40J7jwKG8OgvH8u0qWz/xhfzWdqjkUYutFnsrvTxKs117NJBJl6zJRQv5jPIYIXk2w7fTaaT",
	"aVzJOjTSaSjg6WQ6eQoZOMnLhGmOq2iUfN+82MACkzki9DIOMyuhgFfIP6fS+WEMQu+s8S0730+n8R9l",
	"DaNJDRhvue3/xDOhbKDYxH3P9hTtx4jsDtXgSsnoJwlFH5pGRj3EGfYl7dypIF+irHn539Hsp1dcIQcy",
	"PhHX1grPkoPfsemRVlGc2dnSv3SN37uqdK7WKh3M3/p46e7hib/wVjYuEfbHrwNC6UHyuh1HaN+Nu44s",
	"/jD61isstRfKGoMqFgsksvTASYLZz3JCSouRUEtU/wg0pbPacMtLx9WTziSdjV0YoOglYWRcSGHw3567",
	"hKWdJoQ0Am+152gJa/qczQOfvyStE9Hzj7ZcPwjDbwkrKOCb/JApeRco+fCrtT01PlPA7QfKZ+QQpwDP",
	"zyH1QSn0vgp1vRYqwV8ecC6j4p6NVtzMrGStS9HhLpxc11aWl0jurqMf5IC/DN46VHHFpH1hlQoUGbrQ",
	"BEidaU7F38r2AJuQ/WS4wwr5xp2yFhNi+943TArvUOlKq75LbtZCs0851bPFKxywhZMkG2QkD8Wbh0ae",
	"YCsImTSu4nU6nom5AhmkCQoYWBDOzZEd8XXOxPWjNE5lg+ns8WykHHtQGstd3wv0OL/79KM0SYxtOajW",
	"pJNjfzw0IEblwufIgy+aA/eJ5hM8+wN/0l78pH8pD11hhYRGYSncJ7JT+OyZc2qlfcCMDpbL8mRsjlwa",
	"H483Nu7z2UdMiVHh8FVlwiEKYlE6NaSk36yStShxhbV1DRoWbS1kEKiGApbMrsjzOtYtrefi+fT5FLbX",
	"2/8DAAD//5IL+eIVEgAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	for rawPath, rawFunc := range externalRef0.PathToRawSpec(path.Join(path.Dir(pathToFile), "./types.spec.yaml")) {
		if _, ok := res[rawPath]; ok {
			// it is not possible to compare functions in golang, so always overwrite the old value
		}
		res[rawPath] = rawFunc
	}
	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}

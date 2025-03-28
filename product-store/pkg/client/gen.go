// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	externalRef0 "product-store/pkg/types"

	"github.com/oapi-codegen/runtime"
)

// PutProductCategoryJSONRequestBody defines body for PutProductCategory for application/json ContentType.
type PutProductCategoryJSONRequestBody = externalRef0.ProductCategory

// PutProductJSONRequestBody defines body for PutProduct for application/json ContentType.
type PutProductJSONRequestBody = externalRef0.Product

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Healthz request
	Healthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PutProductCategoryWithBody request with any body
	PutProductCategoryWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PutProductCategory(ctx context.Context, body PutProductCategoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProductCategory request
	GetProductCategory(ctx context.Context, productCategoryName string, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PutProductWithBody request with any body
	PutProductWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PutProduct(ctx context.Context, body PutProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetProduct request
	GetProduct(ctx context.Context, productName string, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) Healthz(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewHealthzRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutProductCategoryWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutProductCategoryRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutProductCategory(ctx context.Context, body PutProductCategoryJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutProductCategoryRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProductCategory(ctx context.Context, productCategoryName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProductCategoryRequest(c.Server, productCategoryName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutProductWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutProductRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PutProduct(ctx context.Context, body PutProductJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPutProductRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) GetProduct(ctx context.Context, productName string, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetProductRequest(c.Server, productName)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewHealthzRequest generates requests for Healthz
func NewHealthzRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/healthz")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPutProductCategoryRequest calls the generic PutProductCategory builder with application/json body
func NewPutProductCategoryRequest(server string, body PutProductCategoryJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPutProductCategoryRequestWithBody(server, "application/json", bodyReader)
}

// NewPutProductCategoryRequestWithBody generates requests for PutProductCategory with any type of body
func NewPutProductCategoryRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/product-categories")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetProductCategoryRequest generates requests for GetProductCategory
func NewGetProductCategoryRequest(server string, productCategoryName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "productCategoryName", runtime.ParamLocationPath, productCategoryName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/product-categories/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPutProductRequest calls the generic PutProduct builder with application/json body
func NewPutProductRequest(server string, body PutProductJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPutProductRequestWithBody(server, "application/json", bodyReader)
}

// NewPutProductRequestWithBody generates requests for PutProduct with any type of body
func NewPutProductRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewGetProductRequest generates requests for GetProduct
func NewGetProductRequest(server string, productName string) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "productName", runtime.ParamLocationPath, productName)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/products/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// HealthzWithResponse request
	HealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthzResponse, error)

	// PutProductCategoryWithBodyWithResponse request with any body
	PutProductCategoryWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutProductCategoryResponse, error)

	PutProductCategoryWithResponse(ctx context.Context, body PutProductCategoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PutProductCategoryResponse, error)

	// GetProductCategoryWithResponse request
	GetProductCategoryWithResponse(ctx context.Context, productCategoryName string, reqEditors ...RequestEditorFn) (*GetProductCategoryResponse, error)

	// PutProductWithBodyWithResponse request with any body
	PutProductWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutProductResponse, error)

	PutProductWithResponse(ctx context.Context, body PutProductJSONRequestBody, reqEditors ...RequestEditorFn) (*PutProductResponse, error)

	// GetProductWithResponse request
	GetProductWithResponse(ctx context.Context, productName string, reqEditors ...RequestEditorFn) (*GetProductResponse, error)
}

type HealthzResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *string
	JSON500      *string
}

// Status returns HTTPResponse.Status
func (r HealthzResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r HealthzResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PutProductCategoryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.ProductCategory
	JSON400      *string
	JSON500      *string
}

// Status returns HTTPResponse.Status
func (r PutProductCategoryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PutProductCategoryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProductCategoryResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.ProductCategory
	JSON404      *string
	JSON500      *string
}

// Status returns HTTPResponse.Status
func (r GetProductCategoryResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProductCategoryResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PutProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.Product
	JSON400      *map[string]interface{}
	JSON404      *string
	JSON500      *string
}

// Status returns HTTPResponse.Status
func (r PutProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PutProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetProductResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *externalRef0.Product
	JSON404      *string
	JSON500      *string
}

// Status returns HTTPResponse.Status
func (r GetProductResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetProductResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// HealthzWithResponse request returning *HealthzResponse
func (c *ClientWithResponses) HealthzWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*HealthzResponse, error) {
	rsp, err := c.Healthz(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseHealthzResponse(rsp)
}

// PutProductCategoryWithBodyWithResponse request with arbitrary body returning *PutProductCategoryResponse
func (c *ClientWithResponses) PutProductCategoryWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutProductCategoryResponse, error) {
	rsp, err := c.PutProductCategoryWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutProductCategoryResponse(rsp)
}

func (c *ClientWithResponses) PutProductCategoryWithResponse(ctx context.Context, body PutProductCategoryJSONRequestBody, reqEditors ...RequestEditorFn) (*PutProductCategoryResponse, error) {
	rsp, err := c.PutProductCategory(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutProductCategoryResponse(rsp)
}

// GetProductCategoryWithResponse request returning *GetProductCategoryResponse
func (c *ClientWithResponses) GetProductCategoryWithResponse(ctx context.Context, productCategoryName string, reqEditors ...RequestEditorFn) (*GetProductCategoryResponse, error) {
	rsp, err := c.GetProductCategory(ctx, productCategoryName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProductCategoryResponse(rsp)
}

// PutProductWithBodyWithResponse request with arbitrary body returning *PutProductResponse
func (c *ClientWithResponses) PutProductWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PutProductResponse, error) {
	rsp, err := c.PutProductWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutProductResponse(rsp)
}

func (c *ClientWithResponses) PutProductWithResponse(ctx context.Context, body PutProductJSONRequestBody, reqEditors ...RequestEditorFn) (*PutProductResponse, error) {
	rsp, err := c.PutProduct(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePutProductResponse(rsp)
}

// GetProductWithResponse request returning *GetProductResponse
func (c *ClientWithResponses) GetProductWithResponse(ctx context.Context, productName string, reqEditors ...RequestEditorFn) (*GetProductResponse, error) {
	rsp, err := c.GetProduct(ctx, productName, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetProductResponse(rsp)
}

// ParseHealthzResponse parses an HTTP response from a HealthzWithResponse call
func ParseHealthzResponse(rsp *http.Response) (*HealthzResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &HealthzResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParsePutProductCategoryResponse parses an HTTP response from a PutProductCategoryWithResponse call
func ParsePutProductCategoryResponse(rsp *http.Response) (*PutProductCategoryResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PutProductCategoryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest externalRef0.ProductCategory
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetProductCategoryResponse parses an HTTP response from a GetProductCategoryWithResponse call
func ParseGetProductCategoryResponse(rsp *http.Response) (*GetProductCategoryResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProductCategoryResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest externalRef0.ProductCategory
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParsePutProductResponse parses an HTTP response from a PutProductWithResponse call
func ParsePutProductResponse(rsp *http.Response) (*PutProductResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PutProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest externalRef0.Product
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseGetProductResponse parses an HTTP response from a GetProductWithResponse call
func ParseGetProductResponse(rsp *http.Response) (*GetProductResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetProductResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest externalRef0.Product
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 404:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON404 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

package api

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type RequestHandler struct {
	ListMux http.Handler
}

func (h *RequestHandler) HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	r, err := convertToHTTPRequest(req)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: 500,
			Body:       "Hello, API Gateway - I broke!",
		}, err
	}
	w := NewAPIGatewayV2ResponseWriter()
	slog.Info("serving request:", "request", r)
	h.ListMux.ServeHTTP(w, r)
	slog.Info("serving response:", "request", w.ToAPIGatewayV2HTTPResponse())
	return w.ToAPIGatewayV2HTTPResponse(), nil
}

// Convert APIGatewayV2HTTPRequest to http.Request
func convertToHTTPRequest(request events.APIGatewayV2HTTPRequest) (*http.Request, error) {
	// Construct the URL
	slog.Info("APIGateway request:", "request", request)
	url := fmt.Sprintf("https://%s%s", request.RequestContext.DomainName, request.RawPath)

	// Create a new http.Request
	httpRequest, err := http.NewRequest(request.RequestContext.HTTP.Method, url, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, values := range request.Headers {
		httpRequest.Header.Add(key, values)
	}

	// Set additional properties as needed

	return httpRequest, nil
}

// APIGatewayV2ResponseWriter is a custom http.ResponseWriter implementation
// that generates an APIGatewayV2HTTPResponse.
type APIGatewayV2ResponseWriter struct {
	StatusCode     int
	Headers        http.Header
	Body           *bytes.Buffer
	MultiValueMode bool
}

// NewAPIGatewayV2ResponseWriter creates a new APIGatewayV2ResponseWriter.
func NewAPIGatewayV2ResponseWriter() *APIGatewayV2ResponseWriter {
	return &APIGatewayV2ResponseWriter{
		Headers: make(http.Header),
		Body:    bytes.NewBuffer(nil),
	}
}

// Header returns the header map.
func (w *APIGatewayV2ResponseWriter) Header() http.Header {
	return w.Headers
}

// Write writes the data to the response body.
func (w *APIGatewayV2ResponseWriter) Write(data []byte) (int, error) {
	return w.Body.Write(data)
}

// WriteHeader sets the status code for the response.
func (w *APIGatewayV2ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

// ToAPIGatewayV2HTTPResponse converts the custom response writer to an APIGatewayV2HTTPResponse.
func (w *APIGatewayV2ResponseWriter) ToAPIGatewayV2HTTPResponse() events.APIGatewayV2HTTPResponse {
	headers := make(map[string]string)
	for key, values := range w.Headers {
		if w.MultiValueMode {
			headers[key] = strings.Join(values, ",")
		} else {
			headers[key] = values[0] // Assuming single value for simplicity
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode:      w.StatusCode,
		Headers:         headers,
		Body:            w.Body.String(),
		IsBase64Encoded: false, // Set to true if the body is base64-encoded
	}
}

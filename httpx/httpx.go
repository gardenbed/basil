// Package httpx provides extensions for http applications and services.
package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Middleware is an interface for wrapping an HTTP handler function.
type Middleware interface {
	Wrap(next http.HandlerFunc) http.HandlerFunc
}

// Error responds to an HTTP request with an error.
// If the given error does not specify a status code, the default status code will be used.
func Error(w http.ResponseWriter, err error, defaultStatusCode int) {
	var e *ServerError
	if errors.As(err, &e) {
		http.Error(w, e.Error(), e.StatusCode())
	} else {
		http.Error(w, err.Error(), defaultStatusCode)
	}
}

// ServerError is a custom error type for errors happening in HTTP handlers.
type ServerError struct {
	error
	statusCode int
}

// NewServerError creates a new HTTP server error.
func NewServerError(err error, statusCode int) *ServerError {
	return &ServerError{err, statusCode}
}

// StatusCode returns the appropriate HTTP status code for the error.
func (e *ServerError) StatusCode() int {
	return e.statusCode
}

// ClientError is a custom error type for errors happening when calling an HTTP endpoint.
type ClientError struct {
	message    string
	statusCode int
}

// NewClientError creates a new HTTP client error.
func NewClientError(resp *http.Response) *ClientError {
	var message string
	if resp.Body != nil {
		respBody := new(struct {
			Error   string `json:"error"`
			Message string `json:"message"`
		})

		if b, err := io.ReadAll(resp.Body); err == nil {
			_ = json.Unmarshal(b, respBody)
			switch {
			case respBody.Error != "":
				message = respBody.Error
			case respBody.Message != "":
				message = respBody.Message
			default:
				message = string(b)
			}
		}
	}

	if message == "" {
		message = fmt.Sprintf("%s %s %d", resp.Request.Method, resp.Request.URL.Path, resp.StatusCode)
	} else {
		message = fmt.Sprintf("%s %s %d: %s", resp.Request.Method, resp.Request.URL.Path, resp.StatusCode, message)
	}

	return &ClientError{
		message:    message,
		statusCode: resp.StatusCode,
	}
}

func (e *ClientError) Error() string {
	return e.message
}

// StatusCode returns the status code of the HTTP response.
func (e *ClientError) StatusCode() int {
	return e.statusCode
}

package httpx

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		defaultStatusCode  int
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "WithError",
			err:                errors.New("invalid credentials"),
			defaultStatusCode:  500,
			expectedStatusCode: 500,
			expectedBody:       "invalid credentials\n",
		},
		{
			name: "WithServerError",
			err: &ServerError{
				errors.New("invalid credentials"),
				403,
			},
			defaultStatusCode:  500,
			expectedStatusCode: 403,
			expectedBody:       "invalid credentials\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			Error(rec, tc.err, tc.defaultStatusCode)

			res := rec.Result()
			b, err := io.ReadAll(res.Body)
			assert.NoError(t, err)
			body := string(b)

			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
			assert.Equal(t, tc.expectedBody, body)
		})
	}
}

func TestServerError(t *testing.T) {
	tests := []struct {
		name               string
		err                error
		statusCode         int
		expectedError      string
		expectedStatusCode int
	}{
		{
			name:               "OK",
			err:                errors.New("invalid credentials"),
			statusCode:         403,
			expectedError:      "invalid credentials",
			expectedStatusCode: 403,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := NewServerError(tc.err, tc.statusCode)

			assert.Equal(t, tc.expectedError, err.Error())
			assert.Equal(t, tc.expectedStatusCode, err.StatusCode())
		})
	}
}

func TestNewClientError(t *testing.T) {
	req, err := http.NewRequest("GET", "/item", nil)
	assert.NoError(t, err)

	tests := []struct {
		name               string
		resp               *http.Response
		expectedError      string
		expectedStatusCode int
	}{
		{
			name: "NilBody",
			resp: &http.Response{
				StatusCode: 500,
				Body:       nil,
				Request:    req,
			},
			expectedError:      "GET /item 500",
			expectedStatusCode: 500,
		},
		{
			name: "EmptyBody",
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("")),
				Request:    req,
			},
			expectedError:      "GET /item 500",
			expectedStatusCode: 500,
		},
		{
			name: "OK_JSONBody_Error",
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader(`{"error":"internal server error"}`)),
				Request:    req,
			},
			expectedError:      "GET /item 500: internal server error",
			expectedStatusCode: 500,
		},
		{
			name: "OK_JSONBody_Message",
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader(`{"message":"internal server error"}`)),
				Request:    req,
			},
			expectedError:      "GET /item 500: internal server error",
			expectedStatusCode: 500,
		},
		{
			name: "OK_TextBody",
			resp: &http.Response{
				StatusCode: 500,
				Body:       io.NopCloser(strings.NewReader("internal server error")),
				Request:    req,
			},
			expectedError:      "GET /item 500: internal server error",
			expectedStatusCode: 500,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := NewClientError(tc.resp)

			assert.Equal(t, tc.expectedError, err.Error())
			assert.Equal(t, tc.expectedStatusCode, err.StatusCode())
		})
	}
}

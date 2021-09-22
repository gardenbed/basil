package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gardenbed/basil/telemetry"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	tests := []struct {
		name                string
		opts                Options
		method              string
		url                 string
		header              http.Header
		next                http.HandlerFunc
		expectedMethod      string
		expectedURL         string
		expectedRoute       string
		expectedStatusCode  int
		expectedStatusClass string
	}{
		{
			name:   "HandlerPanics",
			opts:   Options{},
			method: "GET",
			url:    "/v1/items/00000000-0000-0000-0000-000000000000",
			header: http.Header{},
			next: func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong!")
			},
			expectedStatusCode: 500,
		},
		{
			name:   "Success",
			opts:   Options{},
			method: "GET",
			url:    "/v1/items/00000000-0000-0000-0000-000000000000",
			header: http.Header{
				clientNameHeader: []string{"test-client"},
			},
			next: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/00000000-0000-0000-0000-000000000000",
			expectedRoute:       "/v1/items/:id",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
		{
			name:   "BadRequest",
			opts:   Options{},
			method: "GET",
			url:    "/v1/items/00000000-0000-0000-0000-000000000000",
			header: http.Header{
				clientNameHeader: []string{"test-client"},
			},
			next: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				w.WriteHeader(http.StatusBadRequest)
			},
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/00000000-0000-0000-0000-000000000000",
			expectedRoute:       "/v1/items/:id",
			expectedStatusCode:  400,
			expectedStatusClass: "4xx",
		},
		{
			name:   "InternalServerError",
			opts:   Options{},
			method: "GET",
			url:    "/v1/items/00000000-0000-0000-0000-000000000000",
			header: http.Header{
				clientNameHeader: []string{"test-client"},
			},
			next: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				w.WriteHeader(http.StatusInternalServerError)
			},
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/00000000-0000-0000-0000-000000000000",
			expectedRoute:       "/v1/items/:id",
			expectedStatusCode:  500,
			expectedStatusClass: "5xx",
		},
		{
			name:   "WithRequestMetadata",
			opts:   Options{},
			method: "GET",
			url:    "/v1/items/00000000-0000-0000-0000-000000000000",
			header: http.Header{
				requestUUIDHeader: []string{"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"},
				clientNameHeader:  []string{"test-client"},
			},
			next: func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(2 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
			},
			expectedMethod:      "GET",
			expectedURL:         "/v1/items/00000000-0000-0000-0000-000000000000",
			expectedRoute:       "/v1/items/:id",
			expectedStatusCode:  200,
			expectedStatusClass: "2xx",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			probe := telemetry.NewVoidProbe()
			mid := NewMiddleware(probe, tc.opts)
			assert.NotNil(t, mid)

			// http handler for testing
			handler := mid.Wrap(tc.next)

			// Create an http request
			request := httptest.NewRequest(tc.method, tc.url, nil)
			for k, vals := range tc.header {
				for _, v := range vals {
					request.Header.Add(k, v)
				}
			}

			// Testing
			rec := httptest.NewRecorder()
			handler(rec, request)

			resp := rec.Result()
			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			// TODO: Verify logs
			// TODO: Verify metrics
			// TODO: Verify traces
		})
	}
}

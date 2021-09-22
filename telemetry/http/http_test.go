package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriter(t *testing.T) {
	tests := []struct {
		name        string
		statusCode  int
		statusClass string
	}{
		{"200", 101, "1xx"},
		{"200", 200, "2xx"},
		{"201", 201, "2xx"},
		{"201", 202, "2xx"},
		{"300", 300, "3xx"},
		{"400", 400, "4xx"},
		{"400", 403, "4xx"},
		{"404", 404, "4xx"},
		{"400", 409, "4xx"},
		{"500", 500, "5xx"},
		{"500", 501, "5xx"},
		{"502", 502, "5xx"},
		{"500", 503, "5xx"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			}

			middleware := func(w http.ResponseWriter, r *http.Request) {
				rw := newResponseWriter(w)
				handler(rw, r)

				assert.Equal(t, tc.statusCode, rw.StatusCode)
				assert.Equal(t, tc.statusClass, rw.StatusClass)
			}

			r := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			middleware(w, r)
		})
	}
}

func TestHeaderTextMapCarrier(t *testing.T) {
	tests := []struct {
		name         string
		header       http.Header
		keyValues    [][2]string
		expectedKeys []string
	}{
		{
			name:   "OK",
			header: http.Header{},
			keyValues: [][2]string{
				{"User", "jane"},
				{"Tenant", "testing"},
			},
			expectedKeys: []string{"User", "Tenant"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			carrier := &headerTextMapCarrier{
				Header: tc.header,
			}

			for _, kv := range tc.keyValues {
				carrier.Set(kv[0], kv[1])
				assert.Equal(t, kv[1], carrier.Get(kv[0]))
			}

			assert.Empty(t, carrier.Get("invalid"))

			for _, key := range carrier.Keys() {
				assert.NotEmpty(t, carrier.Get(key))
			}
		})
	}
}

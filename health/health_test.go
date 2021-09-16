package health

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type HealthCheckMock struct {
	InCtx    context.Context
	OutError error
}

type mockChecker struct {
	StringOutString string

	HealthCheckIndex int
	HealthCheckMocks []HealthCheckMock
}

func (m *mockChecker) String() string {
	return m.StringOutString
}

func (m *mockChecker) HealthCheck(ctx context.Context) error {
	i := m.HealthCheckIndex
	m.HealthCheckIndex++
	m.HealthCheckMocks[i].InCtx = ctx
	return m.HealthCheckMocks[i].OutError
}

func TestSetLogger(t *testing.T) {
	l := new(voidLogger)
	SetLogger(l)

	assert.Equal(t, l, logger)
}

func TestRegisterChecker(t *testing.T) {
	c1 := new(mockChecker)
	c2 := new(mockChecker)
	RegisterChecker(c1, c2)

	assert.Equal(t, []Checker{c1, c2}, checkers)
}

func TestSetTimeout(t *testing.T) {
	d := 10 * time.Second
	SetTimeout(d)

	assert.Equal(t, d, timeout)
}

func TestHandlerFunc(t *testing.T) {
	tests := []struct {
		name               string
		checkers           []Checker
		expectedStatusCode int
	}{
		{
			name: "Successful",
			checkers: []Checker{
				&mockChecker{
					HealthCheckMocks: []HealthCheckMock{
						{OutError: nil},
					},
				},
				&mockChecker{
					HealthCheckMocks: []HealthCheckMock{
						{OutError: nil},
					},
				},
			},
			expectedStatusCode: 200,
		},
		{
			name: "Unsuccessful",
			checkers: []Checker{
				&mockChecker{
					HealthCheckMocks: []HealthCheckMock{
						{OutError: nil},
					},
				},
				&mockChecker{
					HealthCheckMocks: []HealthCheckMock{
						{OutError: errors.New("failed to check")},
					},
				},
			},
			expectedStatusCode: 503,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger = new(voidLogger)
			checkers = tc.checkers
			timeout = time.Second

			req := httptest.NewRequest("GET", "/", nil)
			rec := httptest.NewRecorder()
			handler := HandlerFunc()
			handler(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Result().StatusCode)
		})
	}
}

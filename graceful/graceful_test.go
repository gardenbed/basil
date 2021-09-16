package graceful

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type (
	ConnectMock struct {
		OutError error
	}

	DisconnectMock struct {
		InCtx    context.Context
		OutError error
	}

	mockClient struct {
		StringOutString string

		ConnectIndex int
		ConnectMocks []ConnectMock

		DisconnectIndex int
		DisconnectMocks []DisconnectMock
	}
)

func (m *mockClient) String() string {
	return m.StringOutString
}

func (m *mockClient) Connect() error {
	i := m.ConnectIndex
	m.ConnectIndex++
	return m.ConnectMocks[i].OutError
}

func (m *mockClient) Disconnect(ctx context.Context) error {
	i := m.DisconnectIndex
	m.DisconnectIndex++
	m.DisconnectMocks[i].InCtx = ctx
	return m.DisconnectMocks[i].OutError
}

type (
	ListenAndServeMock struct {
		OutError error
	}

	ShutdownMock struct {
		InCtx    context.Context
		OutError error
	}

	mockServer struct {
		StringOutString string

		ListenAndServeIndex int
		ListenAndServeMocks []ListenAndServeMock

		ShutdownIndex int
		ShutdownMocks []ShutdownMock
	}
)

func (m *mockServer) String() string {
	return m.StringOutString
}

func (m *mockServer) ListenAndServe() error {
	i := m.ListenAndServeIndex
	m.ListenAndServeIndex++
	return m.ListenAndServeMocks[i].OutError
}

func (m *mockServer) Shutdown(ctx context.Context) error {
	i := m.ShutdownIndex
	m.ShutdownIndex++
	m.ShutdownMocks[i].InCtx = ctx
	return m.ShutdownMocks[i].OutError
}

func TestSetLogger(t *testing.T) {
	l := new(voidLogger)
	SetLogger(l)

	assert.Equal(t, l, logger)
}

func TestRegisterClient(t *testing.T) {
	c1 := new(mockClient)
	c2 := new(mockClient)
	RegisterClient(c1, c2)

	assert.Equal(t, []Client{c1, c2}, clients)
}

func TestRegisterServer(t *testing.T) {
	s1 := new(mockServer)
	s2 := new(mockServer)
	RegisterServer(s1, s2)

	assert.Equal(t, []Server{s1, s2}, servers)
}

func TestSetMaxRetry(t *testing.T) {
	SetMaxRetry(10)

	assert.Equal(t, 10, maxRetry)
}

func TestSetGracePeriod(t *testing.T) {
	d := 10 * time.Second
	SetGracePeriod(d)

	assert.Equal(t, d, gracePeriod)
}

func TestConnectWithRetry(t *testing.T) {
	tests := []struct {
		name          string
		maxRetry      int
		client        Client
		retries       int
		expectedError error
	}{
		{
			name:     "Successful",
			maxRetry: 2,
			client: &mockClient{
				ConnectMocks: []ConnectMock{
					{OutError: nil},
				},
			},
			retries:       1,
			expectedError: nil,
		},
		{
			name:     "NoRetryLeft",
			maxRetry: 2,
			client: &mockClient{
				ConnectMocks: []ConnectMock{
					{OutError: errors.New("failed to connect")},
				},
			},
			retries:       0,
			expectedError: errors.New("failed to connect"),
		},
		{
			name:     "SuccessfulAfterRetury",
			maxRetry: 2,
			client: &mockClient{
				ConnectMocks: []ConnectMock{
					{OutError: errors.New("failed to connect")},
					{OutError: errors.New("failed to connect")},
					{OutError: nil},
				},
			},
			retries:       2,
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger = new(voidLogger)
			maxRetry = tc.maxRetry

			err := connectWithRetry(tc.client, tc.retries)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestTerminateGracefully(t *testing.T) {
	tests := []struct {
		name          string
		clients       []Client
		servers       []Server
		termServers   bool
		expectedError error
	}{
		{
			name: "ClientError",
			clients: []Client{
				&mockClient{
					DisconnectMocks: []DisconnectMock{
						{OutError: errors.New("failed to disconnect")},
					},
				},
			},
			servers: []Server{
				&mockServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			termServers:   true,
			expectedError: errors.New("failed to disconnect"),
		},
		{
			name: "ServerError",
			clients: []Client{
				&mockClient{
					DisconnectMocks: []DisconnectMock{
						{OutError: nil},
					},
				},
			},
			servers: []Server{
				&mockServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: errors.New("failed to shutdown")},
					},
				},
			},
			termServers:   true,
			expectedError: errors.New("failed to shutdown"),
		},
		{
			name: "ClientAndServerSuccessful",
			clients: []Client{
				&mockClient{
					DisconnectMocks: []DisconnectMock{
						{OutError: nil},
					},
				},
			},
			servers: []Server{
				&mockServer{
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			termServers:   true,
			expectedError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger = new(voidLogger)
			clients = tc.clients
			servers = tc.servers
			maxRetry = 1
			gracePeriod = time.Second

			err := terminateGracefully(tc.termServers)

			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestStartAndWait(t *testing.T) {
	tests := []struct {
		name         string
		clients      []Client
		servers      []Server
		signal       os.Signal
		expectedCode int
	}{
		{
			name: "ClientFailsToConnect",
			clients: []Client{
				&mockClient{
					ConnectMocks: []ConnectMock{
						{OutError: errors.New("failed to connect")},
					},
					DisconnectMocks: []DisconnectMock{
						{OutError: nil},
					},
				},
			},
			servers:      []Server{},
			expectedCode: errorCode,
		},
		{
			name: "ServerFailsToListen",
			clients: []Client{
				&mockClient{
					ConnectMocks: []ConnectMock{
						{OutError: nil},
					},
					DisconnectMocks: []DisconnectMock{
						{OutError: nil},
					},
				},
			},
			servers: []Server{
				&mockServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: errors.New("failed to listen")},
					},
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			expectedCode: errorCode,
		},
		{
			name: "SuccessfulGracefulTermination",
			clients: []Client{
				&mockClient{
					ConnectMocks: []ConnectMock{
						{OutError: nil},
					},
					DisconnectMocks: []DisconnectMock{
						{OutError: nil},
					},
				},
			},
			servers: []Server{
				&mockServer{
					ListenAndServeMocks: []ListenAndServeMock{
						{OutError: nil},
					},
					ShutdownMocks: []ShutdownMock{
						{OutError: nil},
					},
				},
			},
			signal:       syscall.SIGTERM,
			expectedCode: successCode,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger = new(voidLogger)
			clients = tc.clients
			servers = tc.servers
			maxRetry = 0
			gracePeriod = time.Second

			// Get current process
			proc, err := os.FindProcess(os.Getpid())
			assert.NoError(t, err)

			if tc.signal != nil {
				go func() {
					time.Sleep(50 * time.Millisecond)
					_ = proc.Signal(tc.signal)
				}()
			}

			code := StartAndWait()

			assert.Equal(t, tc.expectedCode, code)
		})
	}
}

package grpc

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/stretchr/testify/assert"
)

type (
	SetHeaderMock struct {
		InMD     metadata.MD
		OutError error
	}

	SendHeaderMock struct {
		InMD     metadata.MD
		OutError error
	}

	SetTrailerMock struct {
		InMD metadata.MD
	}

	ContextMock struct {
		OutContext context.Context
	}

	SendMsgMock struct {
		InMsg    interface{}
		OutError error
	}

	RecvMsgMock struct {
		InMsg    interface{}
		OutError error
	}

	MockServerStream struct {
		SetHeaderIndex int
		SetHeaderMocks []SetHeaderMock

		SendHeaderIndex int
		SendHeaderMocks []SendHeaderMock

		SetTrailerIndex int
		SetTrailerMocks []SetTrailerMock

		ContextIndex int
		ContextMocks []ContextMock

		SendMsgIndex int
		SendMsgMocks []SendMsgMock

		RecvMsgIndex int
		RecvMsgMocks []RecvMsgMock
	}
)

func (m *MockServerStream) SetHeader(md metadata.MD) error {
	i := m.SetHeaderIndex
	m.SetHeaderIndex++
	m.SetHeaderMocks[i].InMD = md
	return m.SetHeaderMocks[i].OutError
}

func (m *MockServerStream) SendHeader(md metadata.MD) error {
	i := m.SendHeaderIndex
	m.SendHeaderIndex++
	m.SendHeaderMocks[i].InMD = md
	return m.SendHeaderMocks[i].OutError
}

func (m *MockServerStream) SetTrailer(md metadata.MD) {
	i := m.SetTrailerIndex
	m.SetTrailerIndex++
	m.SetTrailerMocks[i].InMD = md
}

func (m *MockServerStream) Context() context.Context {
	i := m.ContextIndex
	m.ContextIndex++
	return m.ContextMocks[i].OutContext
}

func (m *MockServerStream) SendMsg(msg interface{}) error {
	i := m.SendMsgIndex
	m.SendMsgIndex++
	m.SendMsgMocks[i].InMsg = msg
	return m.SendMsgMocks[i].OutError
}

func (m *MockServerStream) RecvMsg(msg interface{}) error {
	i := m.RecvMsgIndex
	m.RecvMsgIndex++
	m.RecvMsgMocks[i].InMsg = msg
	return m.RecvMsgMocks[i].OutError
}

func TestEndpoint(t *testing.T) {
	tests := []struct {
		name            string
		fullMethod      string
		expectedOK      bool
		expectedPackage string
		expectedService string
		expectedMethod  string
		expectedString  string
	}{
		{
			name:       "Invalid",
			fullMethod: "GetUser",
			expectedOK: false,
		},
		{
			name:            "Valid",
			fullMethod:      "/userPB.UserManager/GetUser",
			expectedOK:      true,
			expectedPackage: "userPB",
			expectedService: "UserManager",
			expectedMethod:  "GetUser",
			expectedString:  "userPB::UserManager::GetUser",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e, ok := parseEndpoint(tc.fullMethod)

			assert.Equal(t, tc.expectedOK, ok)
			assert.Equal(t, tc.expectedPackage, e.Package)
			assert.Equal(t, tc.expectedService, e.Service)
			assert.Equal(t, tc.expectedMethod, e.Method)
			assert.Equal(t, tc.expectedString, e.String())
		})
	}
}

func TestServerStream(t *testing.T) {
	type contextKey string
	baseCtx := context.WithValue(context.Background(), contextKey("key"), "value")
	newCtx := context.WithValue(context.Background(), contextKey("foo"), "bar")

	tests := []struct {
		name        string
		ctx         context.Context
		stream      grpc.ServerStream
		expextedCtx context.Context
	}{
		{
			name: "NoContext",
			ctx:  nil,
			stream: &MockServerStream{
				ContextMocks: []ContextMock{
					{OutContext: baseCtx},
				},
			},
			expextedCtx: baseCtx,
		},
		{
			name: "WithContext",
			ctx:  newCtx,
			stream: &MockServerStream{
				ContextMocks: []ContextMock{
					{OutContext: baseCtx},
				},
			},
			expextedCtx: newCtx,
		},
		{
			name: "AlreadyWrapped",
			ctx:  nil,
			stream: &serverStream{
				ServerStream: &MockServerStream{
					ContextMocks: []ContextMock{
						{OutContext: baseCtx},
					},
				},
			},
			expextedCtx: baseCtx,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ss := ServerStreamWithContext(tc.ctx, tc.stream)

			assert.NotNil(t, ss)
			assert.Equal(t, tc.expextedCtx, ss.Context())
		})
	}
}

func TestMetadataTextMapCarrier(t *testing.T) {
	tests := []struct {
		name         string
		md           *metadata.MD
		keyValues    [][2]string
		expectedKeys []string
	}{
		{
			name: "OK",
			md:   &metadata.MD{},
			keyValues: [][2]string{
				{"user", "jane"},
				{"tenant", "testing"},
			},
			expectedKeys: []string{"user", "tenant"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			carrier := &metadataTextMapCarrier{
				md: tc.md,
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

package grpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/gardenbed/basil/telemetry"

	"github.com/stretchr/testify/assert"
)

func TestClient_UnaryInterceptor(t *testing.T) {
	tests := []struct {
		name             string
		opts             Options
		ctx              context.Context
		method           string
		req              interface{}
		res              interface{}
		cc               *grpc.ClientConn
		callOpts         []grpc.CallOption
		mockInvokerError error
		expectedPackage  string
		expectedService  string
		expectedMethod   string
		expectedStream   bool
		expectedSuccess  bool
	}{
		{
			name:             "InvalidMethod",
			opts:             Options{},
			ctx:              context.Background(),
			method:           "",
			req:              nil,
			res:              nil,
			cc:               nil,
			callOpts:         nil,
			mockInvokerError: nil,
		},
		{
			name: "ExcludedMethods",
			opts: Options{
				ExcludedMethods: []string{"GetItem"},
			},
			ctx:              context.Background(),
			method:           "/itemPB.ItemManager/GetItem",
			req:              nil,
			res:              nil,
			cc:               &grpc.ClientConn{},
			callOpts:         []grpc.CallOption{},
			mockInvokerError: nil,
		},
		{
			name:             "InvokerFails",
			opts:             Options{},
			ctx:              context.Background(),
			method:           "/itemPB.ItemManager/GetItem",
			req:              nil,
			res:              nil,
			cc:               &grpc.ClientConn{},
			callOpts:         []grpc.CallOption{},
			mockInvokerError: errors.New("error on grpc method"),
			expectedPackage:  "itemPB",
			expectedService:  "ItemManager",
			expectedMethod:   "GetItem",
			expectedStream:   false,
			expectedSuccess:  false,
		},
		{
			name:             "InvokerSucceeds",
			opts:             Options{},
			ctx:              context.Background(),
			method:           "/itemPB.ItemManager/GetItem",
			req:              nil,
			res:              nil,
			cc:               &grpc.ClientConn{},
			callOpts:         []grpc.CallOption{},
			mockInvokerError: nil,
			expectedPackage:  "itemPB",
			expectedService:  "ItemManager",
			expectedMethod:   "GetItem",
			expectedStream:   false,
			expectedSuccess:  true,
		},
		{
			name:             "WithRequestUUID",
			opts:             Options{},
			ctx:              telemetry.ContextWithUUID(context.Background(), "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			method:           "/itemPB.ItemManager/GetItem",
			req:              nil,
			res:              nil,
			cc:               &grpc.ClientConn{},
			callOpts:         []grpc.CallOption{},
			mockInvokerError: nil,
			expectedPackage:  "itemPB",
			expectedService:  "ItemManager",
			expectedMethod:   "GetItem",
			expectedStream:   false,
			expectedSuccess:  true,
		},
		{
			name: "WithMetadata",
			opts: Options{},
			ctx: metadata.NewOutgoingContext(context.Background(),
				metadata.New(map[string]string{}),
			),
			method:           "/itemPB.ItemManager/GetItem",
			req:              nil,
			res:              nil,
			cc:               &grpc.ClientConn{},
			callOpts:         []grpc.CallOption{},
			mockInvokerError: nil,
			expectedPackage:  "itemPB",
			expectedService:  "ItemManager",
			expectedMethod:   "GetItem",
			expectedStream:   false,
			expectedSuccess:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			probe := telemetry.NewVoidProbe()
			ci := NewClientInterceptor(probe, tc.opts)
			assert.NotNil(t, ci)

			dialOpts := ci.DialOptions()
			assert.Len(t, dialOpts, 2)

			// grpc invoker for testing
			invoker := func(ctx context.Context, method string, req, res interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				time.Sleep(2 * time.Millisecond)
				return tc.mockInvokerError
			}

			// Testing
			err := ci.unaryInterceptor(tc.ctx, tc.method, tc.req, tc.res, tc.cc, invoker, tc.callOpts...)
			assert.Equal(t, tc.mockInvokerError, err)

			// TODO: Verify logs
			// TODO: Verify metrics
			// TODO: Verify traces
		})
	}
}

func TestClient_StreamInterceptor(t *testing.T) {
	tests := []struct {
		name                 string
		opts                 Options
		ctx                  context.Context
		desc                 *grpc.StreamDesc
		cc                   *grpc.ClientConn
		method               string
		callOpts             []grpc.CallOption
		mockStreamerResponse grpc.ClientStream
		mockStreamerError    error
		expectedPackage      string
		expectedService      string
		expectedMethod       string
		expectedStream       bool
		expectedSuccess      bool
	}{
		{
			name:                 "InvalidMethod",
			opts:                 Options{},
			ctx:                  context.Background(),
			desc:                 nil,
			cc:                   nil,
			method:               "",
			callOpts:             nil,
			mockStreamerResponse: nil,
			mockStreamerError:    nil,
		},
		{
			name: "ExcludedMethods",
			opts: Options{
				ExcludedMethods: []string{"GetItems"},
			},
			ctx:                  context.Background(),
			desc:                 &grpc.StreamDesc{},
			cc:                   &grpc.ClientConn{},
			method:               "/itemPB.ItemManager/GetItems",
			callOpts:             []grpc.CallOption{},
			mockStreamerResponse: nil,
			mockStreamerError:    nil,
		},
		{
			name:                 "StreamerFails",
			opts:                 Options{},
			ctx:                  context.Background(),
			desc:                 &grpc.StreamDesc{},
			cc:                   &grpc.ClientConn{},
			method:               "/itemPB.ItemManager/GetItems",
			callOpts:             []grpc.CallOption{},
			mockStreamerResponse: nil,
			mockStreamerError:    errors.New("error on grpc method"),
			expectedPackage:      "itemPB",
			expectedService:      "ItemManager",
			expectedMethod:       "GetItems",
			expectedStream:       true,
			expectedSuccess:      false,
		},
		{
			name:                 "StreamerSucceeds",
			opts:                 Options{},
			ctx:                  context.Background(),
			desc:                 &grpc.StreamDesc{},
			cc:                   &grpc.ClientConn{},
			method:               "/itemPB.ItemManager/GetItems",
			callOpts:             []grpc.CallOption{},
			mockStreamerResponse: nil,
			mockStreamerError:    nil,
			expectedPackage:      "itemPB",
			expectedService:      "ItemManager",
			expectedMethod:       "GetItems",
			expectedStream:       true,
			expectedSuccess:      true,
		},
		{
			name:                 "WithRequestUUID",
			opts:                 Options{},
			ctx:                  telemetry.ContextWithUUID(context.Background(), "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			desc:                 &grpc.StreamDesc{},
			cc:                   &grpc.ClientConn{},
			method:               "/itemPB.ItemManager/GetItems",
			callOpts:             []grpc.CallOption{},
			mockStreamerResponse: nil,
			mockStreamerError:    nil,
			expectedPackage:      "itemPB",
			expectedService:      "ItemManager",
			expectedMethod:       "GetItems",
			expectedStream:       true,
			expectedSuccess:      true,
		},
		{
			name: "WithMetadata",
			opts: Options{},
			ctx: metadata.NewOutgoingContext(context.Background(),
				metadata.New(map[string]string{}),
			),
			desc:                 &grpc.StreamDesc{},
			cc:                   &grpc.ClientConn{},
			method:               "/itemPB.ItemManager/GetItems",
			callOpts:             []grpc.CallOption{},
			mockStreamerResponse: nil,
			mockStreamerError:    nil,
			expectedPackage:      "itemPB",
			expectedService:      "ItemManager",
			expectedMethod:       "GetItems",
			expectedStream:       true,
			expectedSuccess:      true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			probe := telemetry.NewVoidProbe()
			ci := NewClientInterceptor(probe, tc.opts)
			assert.NotNil(t, ci)

			dialOpts := ci.DialOptions()
			assert.Len(t, dialOpts, 2)

			// grpc streamer for testing
			streamer := func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
				time.Sleep(2 * time.Millisecond)
				return tc.mockStreamerResponse, tc.mockStreamerError
			}

			// Testing
			cs, err := ci.streamInterceptor(tc.ctx, tc.desc, tc.cc, tc.method, streamer, tc.callOpts...)
			assert.Equal(t, tc.mockStreamerResponse, cs)
			assert.Equal(t, tc.mockStreamerError, err)

			// TODO: Verify logs
			// TODO: Verify metrics
			// TODO: Verify traces
		})
	}
}

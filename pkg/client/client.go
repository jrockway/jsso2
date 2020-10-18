package client

import (
	"context"
	"fmt"

	"github.com/jrockway/jsso2/pkg/jssopb"
	"google.golang.org/grpc"
)

// Set is a set of connected JSSO clients.
type Set struct {
	cc         *grpc.ClientConn
	UserClient jssopb.UserClient
}

// Credentials authenticates requests to the JSSO server.
type Credentials struct {
	Token string
}

// GetRequestMetadata implements grpc.PerRPCCredentials.
func (c *Credentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": c.Token,
	}, nil
}

// RequireTransportSecurity implements grpc.PerRPCCredentials.
func (c *Credentials) RequireTransportSecurity() bool {
	// Kind of a bad idea to send your token in the clear, but use your judgement.
	return false
}

// FromCC returns a clientset based on an existing client connection.
func FromCC(cc *grpc.ClientConn) *Set {
	return &Set{
		cc:         cc,
		UserClient: jssopb.NewUserClient(cc),
	}
}

// Dial dials a JSSO server and returns a clientset.
func Dial(ctx context.Context, address, token string, dialopts ...grpc.DialOption) (*Set, error) {
	dialopts = append(dialopts, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&Credentials{Token: token}))

	cc, err := grpc.DialContext(ctx, address, dialopts...)
	if err != nil {
		return nil, fmt.Errorf("dial: %w", err)
	}
	return FromCC(cc), nil
}

// Close closes the clientset's underlying client connection.
func (s *Set) Close() error {
	return s.cc.Close()
}

package client

import (
	"context"
	"fmt"

	"github.com/jrockway/jsso2/pkg/jssopb"
	"google.golang.org/grpc"
)

// Set is a set of connected JSSO clients.
type Set struct {
	cc            *grpc.ClientConn
	UserClient    jssopb.UserClient
	SessionClient jssopb.SessionClient
}

// Credentials authenticates requests to the JSSO server.
type Credentials struct {
	Root   string // Set to authenticate with a root password.
	Token  string // Set to authenticate with a session ID.
	Bearer string // Set to authenticate with a bearer token.
}

// GetRequestMetadata implements grpc.PerRPCCredentials.
func (c *Credentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	result := map[string]string{}
	if c.Token != "" {
		result["Authorization"] = "SessionID " + c.Token
	} else if c.Bearer != "" {
		result["Authorization"] = "Bearer " + c.Bearer
	} else if c.Root != "" {
		result["Authorization"] = "root " + c.Root
	}
	return result, nil
}

// RequireTransportSecurity implements grpc.PerRPCCredentials.
func (c *Credentials) RequireTransportSecurity() bool {
	// Kind of a bad idea to send your token in the clear, but use your judgement.
	return false
}

// FromCC returns a clientset based on an existing client connection.
func FromCC(cc *grpc.ClientConn) *Set {
	return &Set{
		cc:            cc,
		UserClient:    jssopb.NewUserClient(cc),
		SessionClient: jssopb.NewSessionClient(cc),
	}
}

// Dial dials a JSSO server and returns a clientset.
func Dial(ctx context.Context, address string, creds *Credentials, dialopts ...grpc.DialOption) (*Set, error) {
	dialopts = append(dialopts, grpc.WithInsecure(), grpc.WithPerRPCCredentials(creds))

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

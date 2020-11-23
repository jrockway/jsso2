package session

import (
	"context"

	"github.com/jrockway/jsso2/pkg/jssopb"
)

type Service struct {
}

func (s *Service) AuthorizeHTTP(ctx context.Context, req *jssopb.AuthorizeHTTPRequest) (*jssopb.AuthorizeHTTPReply, error) {
	reply := &jssopb.AuthorizeHTTPReply{
		Decision: &jssopb.AuthorizeHTTPReply_Deny{
			Deny: &jssopb.Deny{
				Destination: &jssopb.Deny_Response_{
					Response: &jssopb.Deny_Response{
						ContentType: "text/plain",
						Body:        "Unauthorized.",
					},
				},
			},
		},
	}
	return reply, nil
}

package envoyauthz

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/jrockway/jsso2/pkg/cookies"
	"github.com/jrockway/jsso2/pkg/jssopb"
	protostatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	maxRetries      = 3
	backoffInterval = 10 * time.Millisecond
)

type Config struct {
	Address string `long:"jsso_server_address" env:"JSSO_SERVER_ADDRESS" description:"The URL of JSSO's gRPC server."`
}

type Service struct {
	SessionClient jssopb.SessionClient
}

func (s *Service) Check(ctx context.Context, req *envoy_auth.CheckRequest) (*envoy_auth.CheckResponse, error) {
	reply := &envoy_auth.CheckResponse{
		Status: &protostatus.Status{
			Code: int32(codes.Unavailable),
		},
		HttpResponse: &envoy_auth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_auth.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_InternalServerError,
				},
				Body: "Failed to evaluate authorization decision.",
				Headers: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "content-type",
							Value: "text-plain",
						},
					},
				},
			},
		},
	}

	attrs := req.GetAttributes()
	httpReq := attrs.GetRequest().GetHttp()

	headers := httpReq.GetHeaders()
	if headers == nil {
		headers = make(map[string]string)
	}
	var authorizationHeaders []string
	authorizationHeaders = append(authorizationHeaders, strings.Split(headers["authorization"], ",")...)
	var requestCookies []string
	for _, c := range cookies.Cookies(headers["cookie"]) {
		// I doubt this does well with cookies that contain ",".
		requestCookies = append(requestCookies, c.String())
	}

	requestID := headers["x-request-id"]
	requestURL := &url.URL{
		Scheme: headers["x-forwarded-proto"],
		Host:   httpReq.GetHost(),
		Path:   httpReq.GetPath(),
	}

	authorizeReq := &jssopb.AuthorizeHTTPRequest{
		RequestMethod:        attrs.GetRequest().GetHttp().GetMethod(),
		RequestUri:           requestURL.String(),
		RequestId:            requestID,
		AuthorizationHeaders: authorizationHeaders,
		Cookies:              requestCookies,
		IpAddress:            attrs.GetSource().GetAddress().GetSocketAddress().GetAddress(),
	}
	var errs []string
	for i := 0; i < maxRetries; i++ {
		auth, err := s.SessionClient.AuthorizeHTTP(ctx, authorizeReq)
		if err != nil {
			errs = append(errs, fmt.Sprintf("[call remote AuthorizeHTTP: %v]", err))
			t := time.NewTimer(backoffInterval)
			select {
			case <-t.C:
				continue
			case <-ctx.Done():
				t.Stop()
				var code = codes.Unknown
				switch ctx.Err() {
				case context.DeadlineExceeded:
					code = codes.DeadlineExceeded
				case context.Canceled:
					code = codes.Canceled
				}
				return reply, status.Error(code, fmt.Sprintf("after %d tries: %v", i+1, errs))
			}
		}
		switch decision := auth.Decision.(type) {
		case *jssopb.AuthorizeHTTPReply_Allow:
			allowRes := auth.GetAllow()
			allow := &envoy_auth.OkHttpResponse{}
			reply.Status = &protostatus.Status{
				Code: int32(codes.OK),
			}
			reply.HttpResponse = &envoy_auth.CheckResponse_OkResponse{
				OkResponse: allow,
			}
			for _, h := range allowRes.GetAddHeaders() {
				allow.Headers = append(allow.Headers, &envoy_config_core_v3.HeaderValueOption{
					Append: &wrapperspb.BoolValue{
						Value: false,
					},
					Header: &envoy_config_core_v3.HeaderValue{
						Key:   h.GetKey(),
						Value: h.GetValue(),
					},
				})
			}
			for _, h := range allowRes.GetAppendHeaders() {
				allow.Headers = append(allow.Headers, &envoy_config_core_v3.HeaderValueOption{
					Append: &wrapperspb.BoolValue{
						Value: true,
					},
					Header: &envoy_config_core_v3.HeaderValue{
						Key:   h.GetKey(),
						Value: h.GetValue(),
					},
				})
			}
		case *jssopb.AuthorizeHTTPReply_Deny:
			deny := &envoy_auth.DeniedHttpResponse{}
			reply.Status = &protostatus.Status{
				Code: int32(codes.PermissionDenied),
			}
			reply.HttpResponse = &envoy_auth.CheckResponse_DeniedResponse{
				DeniedResponse: deny,
			}
			switch decision.Deny.GetDestination().(type) {
			case *jssopb.Deny_Redirect_:
				denyRed := decision.Deny.GetRedirect()
				deny.Status = &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_TemporaryRedirect,
				}
				deny.Body = fmt.Sprintf("Not authorized.  Redirecting you to %q", denyRed.GetRedirectUrl())
				deny.Headers = []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "content-type",
							Value: "text/plain",
						},
						Append: &wrapperspb.BoolValue{
							Value: false,
						},
					},
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "location",
							Value: denyRed.GetRedirectUrl(),
						},
						Append: &wrapperspb.BoolValue{
							Value: false,
						},
					},
				}
			case *jssopb.Deny_Response_:
				denyRes := decision.Deny.GetResponse()
				deny.Status = &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_Forbidden,
				}
				deny.Body = denyRes.GetBody()
				deny.Headers = []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "content-type",
							Value: denyRes.GetContentType(),
						},
						Append: &wrapperspb.BoolValue{
							Value: false,
						},
					},
				}
			}
		}
		return reply, nil
	}

	msg := fmt.Sprintf("authorization check failed after %d tries: %v", maxRetries, errs)
	return &envoy_auth.CheckResponse{
		HttpResponse: &envoy_auth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_auth.DeniedHttpResponse{
				Status: &envoy_type_v3.HttpStatus{
					Code: envoy_type_v3.StatusCode_ServiceUnavailable,
				},
				Body: msg,
				Headers: []*envoy_config_core_v3.HeaderValueOption{
					{
						Header: &envoy_config_core_v3.HeaderValue{
							Key:   "content-type",
							Value: "text-plain",
						},
					},
				},
			},
		},
	}, status.Error(codes.Unavailable, msg)
}

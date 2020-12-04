package envoyauthz

import (
	"context"
	"fmt"
	"net/textproto"
	"net/url"
	"strings"
	"time"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type_v3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/sessions"
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
	Address                    string `long:"jsso_server_address" env:"JSSO_SERVER_ADDRESS" description:"The URL of JSSO's gRPC server."`
	AddPlaintextUsernameHeader string `long:"plaintext_username_header" env:"PLAINTEXT_USERNAME_HEADER" description:"If set, send the authenticated user's username in a header with this name."`
}

type Service struct {
	UsernameHeader string
	SessionClient  jssopb.SessionClient
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
	for _, c := range sessions.Cookies(headers["cookie"]) {
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
			headers := map[string][]string{}
			for _, h := range allowRes.GetAddHeaders() {
				k := textproto.CanonicalMIMEHeaderKey(h.GetKey())
				headers[k] = append(headers[k], h.GetValue())
			}
			if _, ok := headers["Cookie"]; !ok {
				allow.HeadersToRemove = append(allow.HeadersToRemove, "cookie")
			}
			if _, ok := headers["Authorization"]; !ok {
				allow.HeadersToRemove = append(allow.HeadersToRemove, "authorization")
			}
			for k, v := range headers {
				// This needs some tweaking.  Envoy is happy to proxy multiple
				// copies of a header, but it doesn't have a way to let us add
				// multiple copies of a header.  We can only append with a , or set
				// a single header.
				//
				// RFC2616 Section 4.2 says: Multiple message-header fields with the
				// same field-name MAY be present in a message if and only if the
				// entire field-value for that header field is defined as a
				// comma-separated list [i.e., #(values)]. It MUST be possible to
				// combine the multiple header fields into one "field-name:
				// field-value" pair, without changing the semantics of the message,
				// by appending each subsequent field-value to the first, each
				// separated by a comma.
				//
				// But I haven't found anything that does that except Envoy when
				// generating a CheckRequest.  Go's http server, for example, treats:
				//
				//   Authorization: foo,bar
				//
				// very differently from:
				//
				//   Authorization: foo
				//   Authorization: bar
				//
				// I suppose this is unlikely to matter in any case that we care
				// about.  Nobody is really sending multiple Authorization headers,
				// and if one of them is for us, we consume that and only set a
				// single Authorization header on the upstream request, so that case
				// works OK.  Cookies we handle specially, because whoever invented
				// Cookies did not care for RFC2616.  RFC7230 at least mentions that
				// (and revises the above text about separators to make it somewhat
				// clear you can't do it in general.)
				joined := strings.Join(v, ",")
				if k == "Cookie" {
					// RFC 6265 4.2.1: ...the user agent will send a Cookie
					// header that conforms to the following grammar:
					// cookie-header = "Cookie:" OWS cookie-string OWS
					// cookie-string = cookie-pair *( ";" SP cookie-pair )
					joined = strings.Join(v, "; ")
				}
				allow.Headers = append(allow.Headers, &envoy_config_core_v3.HeaderValueOption{
					Append: &wrapperspb.BoolValue{
						Value: false,
					},
					Header: &envoy_config_core_v3.HeaderValue{
						Key:   k,
						Value: joined,
					},
				})
			}
			if h := s.UsernameHeader; h != "" {
				allow.Headers = append(allow.Headers, &envoy_config_core_v3.HeaderValueOption{
					Append: &wrapperspb.BoolValue{
						Value: false,
					},
					Header: &envoy_config_core_v3.HeaderValue{
						Key:   h,
						Value: allowRes.GetUsername(),
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

package session

import (
	"context"
	"fmt"
	"net/url"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/jrockway/jsso2/pkg/internalauth"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/redirecttokens"
	"github.com/jrockway/jsso2/pkg/sessions"
	"github.com/jrockway/jsso2/pkg/store"
	"github.com/jrockway/jsso2/pkg/types"
	"github.com/jrockway/jsso2/pkg/web"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	DB          *store.Connection
	Permissions *internalauth.Permissions
	Linker      *web.Linker
	Cookies     *sessions.CookieConfig
	Redirects   *redirecttokens.Config
}

func (s *Service) AuthorizeHTTP(ctx context.Context, req *jssopb.AuthorizeHTTPRequest) (*jssopb.AuthorizeHTTPReply, error) {
	reply := &jssopb.AuthorizeHTTPReply{
		Decision: &jssopb.AuthorizeHTTPReply_Deny{
			Deny: &jssopb.Deny{
				Destination: &jssopb.Deny_Redirect_{
					Redirect: &jssopb.Deny_Redirect{
						RedirectUrl: s.Linker.LoginPage(),
					},
				},
			},
		},
	}
	l := ctxzap.Extract(ctx)
	// Check that the request URL is valid
	parsedURL, err := url.Parse(req.GetRequestUri())
	if err != nil {
		return reply, status.Error(codes.InvalidArgument, fmt.Errorf("parse request uri: %w", err).Error())
	}

	if redirectToken, err := s.Redirects.New(parsedURL.String()); err != nil {
		l.Warn("could not mint redirect token", zap.String("url", parsedURL.String()), zap.Error(err))
		reply.GetDeny().GetRedirect().RedirectUrl = s.Linker.LoginPage()
	} else {
		reply.GetDeny().GetRedirect().RedirectUrl = s.Linker.LoginPageWithRedirect(redirectToken)
	}

	// Check is the proxy's user is allowed to perform this check.
	if err := s.Permissions.AllowAuthorizeHTTP(ctx, sessions.MustFromContext(ctx).GetUser()); err != nil {
		return reply, err
	}

	// Extract the end user's session from the request.
	ss, unusedAuth, unusedCookies := s.Cookies.SessionsFromAny(req.GetAuthorizationHeaders(), req.GetCookies())
	session, errs := s.DB.AuthenticateUser(ctx, l, ss, unusedAuth, unusedCookies)
	if session == nil {
		switch len(errs) {
		case 0:
			reply.GetDeny().Reason = "no authentication material provided"
		case 1:
			reply.GetDeny().Reason = fmt.Sprintf("%v", errs[0])
		default:
			reply.GetDeny().Reason = fmt.Sprintf("%d errors: %v", len(errs), errs)
			return reply, nil
		}
	}

	// Check that the access control policy allows this user to visit the target website.
	if err := s.Permissions.AllowWebVisit(ctx, session, parsedURL); err != nil {
		reply.GetDeny().Reason = err.Error()
		return reply, nil
	}

	allow := &jssopb.Allow{
		Username: session.GetUser().GetUsername(),
	}
	for _, u := range unusedAuth {
		// We can't really tell which authorization headers were intended for us if they are
		// malformed, so pass all authorization headers along.
		allow.AddHeaders = append(allow.AddHeaders, &types.Header{
			Key:   "Authorization",
			Value: u.Value,
		})
	}
	for _, u := range unusedCookies {
		if u.Err != nil {
			allow.AddHeaders = append(allow.AddHeaders, &types.Header{
				Key:   "Cookie",
				Value: u.Cookie.String(),
			})
		}
	}
	reply = &jssopb.AuthorizeHTTPReply{
		Decision: &jssopb.AuthorizeHTTPReply_Allow{
			Allow: allow,
		},
	}
	return reply, nil
}

package testserver

import (
	"testing"

	"github.com/jrockway/jsso2/pkg/jsso/enrollment"
	"github.com/jrockway/jsso2/pkg/jsso/login"
	"github.com/jrockway/jsso2/pkg/jsso/user"
	"github.com/jrockway/jsso2/pkg/jssopb"
	"github.com/jrockway/jsso2/pkg/jtesting"
	"github.com/jrockway/jsso2/pkg/store"
	"google.golang.org/grpc"
)

func Setup(t *testing.T, e *jtesting.E, s *grpc.Server) {
	db := store.MustGetTestDB(t, e)
	jssopb.RegisterEnrollmentService(s, jssopb.NewEnrollmentService(&enrollment.Service{}))
	jssopb.RegisterUserService(s, jssopb.NewUserService(&user.Service{DB: db}))
	jssopb.RegisterLoginService(s, jssopb.NewLoginService(&login.Service{}))
}

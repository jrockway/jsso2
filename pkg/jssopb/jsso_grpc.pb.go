// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package jssopb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserClient interface {
	// Add a new user.  The user won't be able to log in until they visit the
	// enrollment URL.
	Add(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserReply, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

var userAddStreamDesc = &grpc.StreamDesc{
	StreamName: "Add",
}

func (c *userClient) Add(ctx context.Context, in *AddUserRequest, opts ...grpc.CallOption) (*AddUserReply, error) {
	out := new(AddUserReply)
	err := c.cc.Invoke(ctx, "/jsso.User/Add", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserService is the service API for User service.
// Fields should be assigned to their respective handler implementations only before
// RegisterUserService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type UserService struct {
	// Add a new user.  The user won't be able to log in until they visit the
	// enrollment URL.
	Add func(context.Context, *AddUserRequest) (*AddUserReply, error)
}

func (s *UserService) add(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Add(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/jsso.User/Add",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Add(ctx, req.(*AddUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterUserService registers a service implementation with a gRPC server.
func RegisterUserService(s grpc.ServiceRegistrar, srv *UserService) {
	srvCopy := *srv
	if srvCopy.Add == nil {
		srvCopy.Add = func(context.Context, *AddUserRequest) (*AddUserReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Add not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "jsso.User",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Add",
				Handler:    srvCopy.add,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "jsso.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewUserService creates a new UserService containing the
// implemented methods of the User service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewUserService(s interface{}) *UserService {
	ns := &UserService{}
	if h, ok := s.(interface {
		Add(context.Context, *AddUserRequest) (*AddUserReply, error)
	}); ok {
		ns.Add = h.Add
	}
	return ns
}

// UnstableUserService is the service API for User service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableUserService interface {
	// Add a new user.  The user won't be able to log in until they visit the
	// enrollment URL.
	Add(context.Context, *AddUserRequest) (*AddUserReply, error)
}

// LoginClient is the client API for Login service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoginClient interface {
	Start(ctx context.Context, in *StartLoginRequest, opts ...grpc.CallOption) (*StartLoginReply, error)
}

type loginClient struct {
	cc grpc.ClientConnInterface
}

func NewLoginClient(cc grpc.ClientConnInterface) LoginClient {
	return &loginClient{cc}
}

var loginStartStreamDesc = &grpc.StreamDesc{
	StreamName: "Start",
}

func (c *loginClient) Start(ctx context.Context, in *StartLoginRequest, opts ...grpc.CallOption) (*StartLoginReply, error) {
	out := new(StartLoginReply)
	err := c.cc.Invoke(ctx, "/jsso.Login/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoginService is the service API for Login service.
// Fields should be assigned to their respective handler implementations only before
// RegisterLoginService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type LoginService struct {
	Start func(context.Context, *StartLoginRequest) (*StartLoginReply, error)
}

func (s *LoginService) start(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/jsso.Login/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Start(ctx, req.(*StartLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterLoginService registers a service implementation with a gRPC server.
func RegisterLoginService(s grpc.ServiceRegistrar, srv *LoginService) {
	srvCopy := *srv
	if srvCopy.Start == nil {
		srvCopy.Start = func(context.Context, *StartLoginRequest) (*StartLoginReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "jsso.Login",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Start",
				Handler:    srvCopy.start,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "jsso.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewLoginService creates a new LoginService containing the
// implemented methods of the Login service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewLoginService(s interface{}) *LoginService {
	ns := &LoginService{}
	if h, ok := s.(interface {
		Start(context.Context, *StartLoginRequest) (*StartLoginReply, error)
	}); ok {
		ns.Start = h.Start
	}
	return ns
}

// UnstableLoginService is the service API for Login service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableLoginService interface {
	Start(context.Context, *StartLoginRequest) (*StartLoginReply, error)
}

// EnrollmentClient is the client API for Enrollment service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EnrollmentClient interface {
	Start(ctx context.Context, in *StartEnrollmentRequest, opts ...grpc.CallOption) (*StartEnrollmentReply, error)
}

type enrollmentClient struct {
	cc grpc.ClientConnInterface
}

func NewEnrollmentClient(cc grpc.ClientConnInterface) EnrollmentClient {
	return &enrollmentClient{cc}
}

var enrollmentStartStreamDesc = &grpc.StreamDesc{
	StreamName: "Start",
}

func (c *enrollmentClient) Start(ctx context.Context, in *StartEnrollmentRequest, opts ...grpc.CallOption) (*StartEnrollmentReply, error) {
	out := new(StartEnrollmentReply)
	err := c.cc.Invoke(ctx, "/jsso.Enrollment/Start", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EnrollmentService is the service API for Enrollment service.
// Fields should be assigned to their respective handler implementations only before
// RegisterEnrollmentService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type EnrollmentService struct {
	Start func(context.Context, *StartEnrollmentRequest) (*StartEnrollmentReply, error)
}

func (s *EnrollmentService) start(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartEnrollmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Start(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/jsso.Enrollment/Start",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Start(ctx, req.(*StartEnrollmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterEnrollmentService registers a service implementation with a gRPC server.
func RegisterEnrollmentService(s grpc.ServiceRegistrar, srv *EnrollmentService) {
	srvCopy := *srv
	if srvCopy.Start == nil {
		srvCopy.Start = func(context.Context, *StartEnrollmentRequest) (*StartEnrollmentReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Start not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "jsso.Enrollment",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Start",
				Handler:    srvCopy.start,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "jsso.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewEnrollmentService creates a new EnrollmentService containing the
// implemented methods of the Enrollment service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewEnrollmentService(s interface{}) *EnrollmentService {
	ns := &EnrollmentService{}
	if h, ok := s.(interface {
		Start(context.Context, *StartEnrollmentRequest) (*StartEnrollmentReply, error)
	}); ok {
		ns.Start = h.Start
	}
	return ns
}

// UnstableEnrollmentService is the service API for Enrollment service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableEnrollmentService interface {
	Start(context.Context, *StartEnrollmentRequest) (*StartEnrollmentReply, error)
}

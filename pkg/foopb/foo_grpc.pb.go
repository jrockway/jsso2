// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package foopb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// NameServiceClient is the client API for NameService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NameServiceClient interface {
	// Transform a name.
	TransformName(ctx context.Context, in *TransformNameRequest, opts ...grpc.CallOption) (*TransformNameReply, error)
}

type nameServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNameServiceClient(cc grpc.ClientConnInterface) NameServiceClient {
	return &nameServiceClient{cc}
}

var nameServiceTransformNameStreamDesc = &grpc.StreamDesc{
	StreamName: "TransformName",
}

func (c *nameServiceClient) TransformName(ctx context.Context, in *TransformNameRequest, opts ...grpc.CallOption) (*TransformNameReply, error) {
	out := new(TransformNameReply)
	err := c.cc.Invoke(ctx, "/foo.NameService/TransformName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NameServiceService is the service API for NameService service.
// Fields should be assigned to their respective handler implementations only before
// RegisterNameServiceService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type NameServiceService struct {
	// Transform a name.
	TransformName func(context.Context, *TransformNameRequest) (*TransformNameReply, error)
}

func (s *NameServiceService) transformName(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TransformNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.TransformName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/foo.NameService/TransformName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.TransformName(ctx, req.(*TransformNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterNameServiceService registers a service implementation with a gRPC server.
func RegisterNameServiceService(s grpc.ServiceRegistrar, srv *NameServiceService) {
	srvCopy := *srv
	if srvCopy.TransformName == nil {
		srvCopy.TransformName = func(context.Context, *TransformNameRequest) (*TransformNameReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method TransformName not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "foo.NameService",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "TransformName",
				Handler:    srvCopy.transformName,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "foo.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewNameServiceService creates a new NameServiceService containing the
// implemented methods of the NameService service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewNameServiceService(s interface{}) *NameServiceService {
	ns := &NameServiceService{}
	if h, ok := s.(interface {
		TransformName(context.Context, *TransformNameRequest) (*TransformNameReply, error)
	}); ok {
		ns.TransformName = h.TransformName
	}
	return ns
}

// UnstableNameServiceService is the service API for NameService service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableNameServiceService interface {
	// Transform a name.
	TransformName(context.Context, *TransformNameRequest) (*TransformNameReply, error)
}
package authorization_server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type AuthorizationServerDesc struct {
	svc AuthorizationServerServer
}

func NewAuthorizationServerDesc(svc AuthorizationServerServer) *AuthorizationServerDesc {
	return &AuthorizationServerDesc{
		svc: svc,
	}
}

func (d *AuthorizationServerDesc) RegisterGRPC(s *grpc.Server) {
	RegisterAuthorizationServerServer(s, d.svc)
}

func (d *AuthorizationServerDesc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterAuthorizationServerHandlerServer(ctx, mux, d.svc)
}

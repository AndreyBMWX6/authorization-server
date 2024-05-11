package authentication_server

import (
	"context"
	_ "embed"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

//go:embed authentication_server.swagger.json
var swaggerJSON []byte

type AuthorizationServerDesc struct {
	svc AuthenticationServerServer
}

func NewAuthorizationServerDesc(svc AuthenticationServerServer) *AuthorizationServerDesc {
	return &AuthorizationServerDesc{
		svc: svc,
	}
}

func (d *AuthorizationServerDesc) RegisterGRPC(s *grpc.Server) {
	RegisterAuthenticationServerServer(s, d.svc)
}

func (d *AuthorizationServerDesc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux) error {
	return RegisterAuthenticationServerHandlerServer(ctx, mux, d.svc)
}

func (d *AuthorizationServerDesc) SwaggerDef() []byte {
	return swaggerJSON
}

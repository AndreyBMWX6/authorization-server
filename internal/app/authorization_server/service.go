package authorization_server

import (
	"authorization-server/internal/app"
	desc "authorization-server/pkg/api/authorization_server"
)

type Implementation struct {
	desc.UnimplementedAuthorizationServerServer
}

func NewAuthorizationServer() *Implementation {
	return &Implementation{}
}

func (i *Implementation) GetDescription() app.ServiceDesc {
	return desc.NewAuthorizationServerDesc(i)
}

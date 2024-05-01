package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	httpPort = 7000
	grpcPort = 7002
)

type Service interface {
	GetDescription() ServiceDesc
}

type ServiceDesc interface {
	RegisterGRPC(*grpc.Server)
	RegisterGateway(context.Context, *runtime.ServeMux) error
	//todo: swagger UI + interceptors
	//SwaggerDef() []byte
	//WithHTTPUnaryInterceptor(grpc.UnaryServerInterceptor)
}

type App struct {
	lis        *listeners
	grpcServer *grpc.Server
	// todo: переделать runtime.mux на chi.Router
	httpServer *runtime.ServeMux
	desc       ServiceDesc
	//todo: добавить порты
}

type listeners struct {
	grpc net.Listener
}

func NewApp() *App {
	return &App{
		lis: &listeners{},
	}
}

func (a *App) Run(ctx context.Context, service Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	desc := service.GetDescription()
	a.desc = desc

	a.initGRPC()
	a.initHTTP()
	a.runGRPC()
	a.runHTTP()

	//todo: add closer
	for {
	}

	return nil
}

func (a *App) initGRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	a.lis.grpc = lis

	a.grpcServer = grpc.NewServer()
	a.desc.RegisterGRPC(a.grpcServer)
	// включаем grpc рефлексию
	reflection.Register(a.grpcServer)
}

func (a *App) initHTTP() {
	a.httpServer = runtime.NewServeMux()
	err := a.desc.RegisterGateway(context.Background(), a.httpServer)
	if err != nil {
		log.Fatalf("failed to refister gateway")
	}
	//todo:добавить swagger UI
}

func (a *App) runGRPC() {
	go func() {
		if err := a.grpcServer.Serve(a.lis.grpc); err != nil {
			log.Fatalf("failed to serve grpc server: %s", err.Error())
		}
	}()
}

func (a *App) runHTTP() {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("localhost:%d", httpPort), a.httpServer); err != nil {

		}
	}()
}

package scratch

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	httpPort      = 7000
	adminHttpPort = 7001
	grpcPort      = 7002
)

type Service interface {
	GetDescription() ServiceDesc
}

type ServiceDesc interface {
	RegisterGRPC(*grpc.Server)
	RegisterGateway(context.Context, *runtime.ServeMux) error
	SwaggerDef() []byte
	//todo: if needed, add interceptors
	//WithHTTPUnaryInterceptor(grpc.UnaryServerInterceptor)
}

type App struct {
	lis          *listeners
	grpcServer   *grpc.Server
	publicServer *chi.Mux
	adminServer  *chi.Mux
	desc         ServiceDesc
	//todo: добавить порты
}

type listeners struct {
	grpc      net.Listener
	http      net.Listener
	httpAdmin net.Listener
}

func New() (*App, error) {
	a := &App{}

	lis, err := newListeners()
	if err != nil {
		return nil, errors.Wrap(err, "create listeners")
	}
	a.lis = lis
	a.initHTTP()

	return a, nil
}

func newListeners() (*listeners, error) {
	//todo: add host option
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		return nil, errors.Wrap(err, "listen grpc port")
	}
	httpLis, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		return nil, errors.Wrap(err, "listen public http port")
	}
	adminHttpLis, err := net.Listen("tcp", fmt.Sprintf(":%d", adminHttpPort))
	if err != nil {
		return nil, errors.Wrap(err, "listen admin http port")
	}

	return &listeners{
		grpc:      grpcLis,
		http:      httpLis,
		httpAdmin: adminHttpLis,
	}, nil
}

func (a *App) initGRPC() {
	a.grpcServer = grpc.NewServer()
	a.desc.RegisterGRPC(a.grpcServer)
	// включаем grpc рефлексию
	reflection.Register(a.grpcServer)
}

func (a *App) initHTTP() {
	a.initAdminHTTP()
	a.initPublicHTTP()
}

func (a *App) initPublicHTTP() {
	a.publicServer = chi.NewMux()

	// setting cors policy
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodHead,
		},
	})

	a.publicServer.Use(c.Handler)
}

func (a *App) initPublicHttpHandlers(desc ServiceDesc) {
	if desc != nil {
		mux, err := NewGatewayMux(desc, a.publicServer)
		if err != nil {
			log.Fatalf("failed to init gateway: %s", err.Error())
		}

		if err := desc.RegisterGateway(context.Background(), mux); err != nil {
			log.Fatalf("failed to register gateway")
		}
	}
}

func (a *App) initAdminHTTP() {
	a.adminServer = chi.NewMux()
	// SwaggerUI
	a.initSwaggerHandlers()
}

func (a *App) initSwaggerHandlers() {
	publicPort := strconv.Itoa(a.lis.http.Addr().(*net.TCPAddr).Port)
	a.adminServer.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		if h, _, err := net.SplitHostPort(r.Host); err == nil {
			r.Host = net.JoinHostPort(h, publicPort)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(swaggerDef(a.desc, r.Host))
	})

	a.adminServer.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/", http.StatusMovedPermanently)
	})

	a.adminServer.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
}

func NewGatewayMux(desc ServiceDesc, cr chi.Router) (*runtime.ServeMux, error) {
	mux := runtime.NewServeMux()
	if err := mountHandlersFromSwagger(desc, cr, mux); err != nil {
		return nil, errors.Wrap(err, "mount swagger handlers")
	}

	return mux, nil
}

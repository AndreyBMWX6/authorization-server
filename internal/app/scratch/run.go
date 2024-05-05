package scratch

import (
	"context"
	"log"
	"net/http"
)

func (a *App) Run(ctx context.Context, service Service) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	desc := service.GetDescription()
	a.desc = desc

	a.initGRPC()
	a.initPublicHttpHandlers(desc)

	a.runGRPC()
	a.runHTTP()

	//todo: add closer
	for {
	}

	return nil
}

func (a *App) runGRPC() {
	go func() {
		if err := a.grpcServer.Serve(a.lis.grpc); err != nil {
			log.Fatalf("grpc server: %s", err.Error())
		}
	}()
}

func (a *App) runHTTP() {
	publicServer := http.Server{Handler: a.publicServer}
	adminServer := http.Server{Handler: a.adminServer}
	go func() {
		if err := publicServer.Serve(a.lis.http); err != nil {
			log.Fatalf("http public server: %s", err.Error())
		}
	}()
	go func() {
		if err := adminServer.Serve(a.lis.httpAdmin); err != nil {
			log.Fatalf("http admin server: %s", err.Error())
		}
	}()
}

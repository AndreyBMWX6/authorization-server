package scratch

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/analysis"
	"github.com/go-openapi/spec"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
)

func swaggerDef(desc ServiceDesc, host string) []byte {
	var s spec.Swagger
	def := desc.SwaggerDef()

	if err := json.Unmarshal(def, &s); err != nil {
		//todo: logger with err levels
		log.Printf("failed to unmarshal swaggerDef: %v", err)
		return def
	}

	s.Host = host

	var err error
	if def, err = json.Marshal(s); err != nil {
		//todo: logger with err levels
		log.Printf("failed to marshal swaggerDef: %v", err)
	}

	return def
}

func mountHandlersFromSwagger(desc ServiceDesc, chiMux chi.Router, h *runtime.ServeMux) error {
	doc := &spec.Swagger{}

	if err := doc.UnmarshalJSON(desc.SwaggerDef()); err != nil {
		return errors.Wrap(err, "unmarshal swaggerDef")
	}

	swag := analysis.New(doc)
	for k, v := range swag.AllPaths() {
		if v.Get != nil {
			chiMux.Method(http.MethodGet, k, h)
		}
		if v.Post != nil {
			chiMux.Method(http.MethodPost, k, h)
		}
		if v.Put != nil {
			chiMux.Method(http.MethodPut, k, h)
		}
		if v.Delete != nil {
			chiMux.Method(http.MethodDelete, k, h)
		}
		if v.Patch != nil {
			chiMux.Method(http.MethodPatch, k, h)
		}
		if v.Options != nil {
			chiMux.Method(http.MethodOptions, k, h)
		}
		if v.Head != nil {
			chiMux.Method(http.MethodHead, k, h)
		}
	}
	return nil
}

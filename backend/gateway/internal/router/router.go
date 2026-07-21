package router

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/sekaishopml/cytaxi/backend/gateway/internal/middleware"
)

type GatewayRouter struct {
	mux      *http.ServeMux
	backends map[string]string
	mw       *middleware.Chain
	logger   *slog.Logger
}

func New(backends map[string]string, mw *middleware.Chain, logger *slog.Logger) *GatewayRouter {
	return &GatewayRouter{
		mux:      http.NewServeMux(),
		backends: backends,
		mw:       mw,
		logger:   logger,
	}
}

func (r *GatewayRouter) Handler() http.Handler {
	return r.mw.Apply(r.mux)
}

func (r *GatewayRouter) Health(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"api-gateway","version":"1.0"}`))
}

func (r *GatewayRouter) RegisterRoute(service, method, path, targetService string) {
	backendURL, ok := r.backends[targetService]
	if !ok {
		r.logger.Warn("backend not configured", "service", targetService)
		return
	}

	target, _ := url.Parse(backendURL)
	proxy := httputil.NewSingleHostReverseProxy(target)

		// Calcular el prefijo a remover a partir del path registrado.
		// El path viene como "/api/v1/<engine>/{path...}" o "/api/v1/<engine>s/{path...}".
		// El prefijo a remover es la parte fija antes de "/{path...}".
		prefixToStrip := strings.Split(path, "/{path...")[0]
		if prefixToStrip == "" {
			prefixToStrip = "/api/v1/" + service
		}
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Remover prefijo /api/v1/<engine> del path
		newPath := strings.TrimPrefix(req.URL.Path, prefixToStrip)
		if !strings.HasPrefix(newPath, "/") {
			newPath = "/" + newPath
		}
		req.URL.Path = newPath
	}

	pattern := method + " " + path
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		r.logger.Info("gateway request",
			"method", req.Method,
			"path", req.URL.Path,
			"target", targetService,
		)
		proxy.ServeHTTP(w, req)
	})
}

func (r *GatewayRouter) RegisterRoutes() {
	engines := map[string]string{
		"customer":     "customer",
		"driver":       "driver",
		"trip":         "trip",
		"pricing":      "pricing",
		"payment":      "payment",
		"notification": "notification",
		"admin":        "admin",
		"analytics":    "analytics",
		"matching":     "matching",
		"geo":          "geo",
	}

	for engine, backend := range engines {
		// Wildcard pattern para capturar sub-rutas (Go 1.22+)
		// Registramos tanto la forma singular (/api/v1/<engine>/...) como la
		// plural (/api/v1/<engine>s/...) porque los frontends usan ambas.
		r.RegisterRoute(engine, "GET", "/api/v1/"+engine+"/{path...}", backend)
		r.RegisterRoute(engine, "POST", "/api/v1/"+engine+"/{path...}", backend)
		r.RegisterRoute(engine, "PUT", "/api/v1/"+engine+"/{path...}", backend)
		r.RegisterRoute(engine, "DELETE", "/api/v1/"+engine+"/{path...}", backend)
		r.RegisterRoute(engine, "GET", "/api/v1/"+engine+"/health", backend)

		plural := engine + "s"
		if plural != engine {
			r.RegisterRoute(engine, "GET", "/api/v1/"+plural+"/{path...}", backend)
			r.RegisterRoute(engine, "POST", "/api/v1/"+plural+"/{path...}", backend)
			r.RegisterRoute(engine, "PUT", "/api/v1/"+plural+"/{path...}", backend)
			r.RegisterRoute(engine, "DELETE", "/api/v1/"+plural+"/{path...}", backend)
			r.RegisterRoute(engine, "GET", "/api/v1/"+plural+"/health", backend)
		}
	}

	r.mux.HandleFunc("GET /health", r.Health)
	r.mux.HandleFunc("GET /api/v1/health", r.Health)
}

func (r *GatewayRouter) GetRoute(ctx context.Context, path string) (string, bool) {
	return "", false
}

package router

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

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
	}

	for engine, backend := range engines {
		r.RegisterRoute(engine, "GET", "/api/v1/"+engine+"/", backend)
		r.RegisterRoute(engine, "POST", "/api/v1/"+engine+"/", backend)
		r.RegisterRoute(engine, "GET", "/api/v1/"+engine+"/health", backend)
	}

	r.mux.HandleFunc("GET /health", r.Health)
	r.mux.HandleFunc("GET /api/v1/health", r.Health)
}

func (r *GatewayRouter) GetRoute(ctx context.Context, path string) (string, bool) {
	return "", false
}

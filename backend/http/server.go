package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	TLS          *tls.Config
}

func DefaultConfig() Config {
	return Config{
		Host:         "0.0.0.0",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

type server struct {
	inner  *http.Server
	config Config
}

func New(handler http.Handler, cfg Config) Server {
	addr := net.JoinHostPort(cfg.Host, fmt.Sprintf("%d", cfg.Port))

	return &server{
		config: cfg,
		inner: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
			TLSConfig:    cfg.TLS,
		},
	}
}

func (s *server) Start() error {
	ln, err := net.Listen("tcp", s.inner.Addr)
	if err != nil {
		return fmt.Errorf("http server: listen: %w", err)
	}

	if s.config.TLS != nil {
		return s.inner.ServeTLS(ln, "", "")
	}
	return s.inner.Serve(ln)
}

func (s *server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.inner.Shutdown(ctx)
}

func (s *server) Route(method, path string, handler http.HandlerFunc) {
	mux, ok := s.inner.Handler.(*http.ServeMux)
	if !ok {
		return
	}
	mux.HandleFunc(fmt.Sprintf("%s %s", method, path), handler)
}

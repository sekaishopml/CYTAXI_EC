package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/application"
	"github.com/sekaishopml/cytaxi/backend/engines/policy/internal/policy/config"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	loader := application.NewMemoryLoader(nil)
	registry := application.NewPolicyRegistry(loader)
	condEval := application.NewConditionEvaluator()
	evaluator := application.NewRuleEvaluator(condEval)
	engine := application.NewPolicyEngine(registry, evaluator, logger)
	_ = engine

	logger.Info("policy engine starting", "port", cfg.Port, "policies_dir", cfg.PoliciesDir)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok","service":"policy-engine"}`))
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down")
	srv.Shutdown(context.Background())
}

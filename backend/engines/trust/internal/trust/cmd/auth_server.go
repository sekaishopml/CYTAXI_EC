package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/infrastructure/auth"
	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/infrastructure/providers"
)

type AuthServer struct {
	tokenManager  *auth.Manager
	roleManager   *auth.RoleManager
	providerReg   *providers.Registry
	logger        *slog.Logger
}

func NewAuthServer(tokenManager *auth.Manager, roleManager *auth.RoleManager, providerReg *providers.Registry, logger *slog.Logger) *AuthServer {
	return &AuthServer{
		tokenManager: tokenManager,
		roleManager:  roleManager,
		providerReg:  providerReg,
		logger:       logger,
	}
}

func (s *AuthServer) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		Method   string `json:"method"`
		Name     string `json:"name"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	method := providers.AuthMethod(req.Method)
	if method == "" { method = providers.MethodEmailPassword }

	p, err := s.providerReg.Get(method)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	result, err := p.Register(r.Context(), providers.Credential{
		Email: req.Email, Password: req.Password, Phone: req.Phone,
	})
	if err != nil || !result.Success {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": result.Error.Error()})
		return
	}

	result.Identity.Name = req.Name

	tokens, err := s.tokenManager.GenerateTokens(auth.Claims{
		Sub:   result.Identity.ID,
		Email: result.Identity.Email,
		Phone: result.Identity.Phone,
		Name:  result.Identity.Name,
		Role:  result.Identity.Role,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token generation failed"})
		return
	}

	s.logger.Info("user registered", "user_id", result.Identity.ID, "method", method)
	writeJSON(w, http.StatusCreated, map[string]any{
		"status": "registered",
		"user":   result.Identity,
		"tokens": tokens,
	})
}

func (s *AuthServer) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
		OTP      string `json:"otp"`
		Method   string `json:"method"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	method := providers.AuthMethod(req.Method)
	if method == "" { method = providers.MethodEmailPassword }

	p, err := s.providerReg.Get(method)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	result, err := p.Authenticate(r.Context(), providers.Credential{
		Email: req.Email, Password: req.Password, Phone: req.Phone, OTP: req.OTP,
	})
	if err != nil || !result.Success {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
		return
	}

	tokens, err := s.tokenManager.GenerateTokens(auth.Claims{
		Sub:   result.Identity.ID,
		Email: result.Identity.Email,
		Phone: result.Identity.Phone,
		Name:  result.Identity.Name,
		Role:  result.Identity.Role,
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "token generation failed"})
		return
	}

	s.logger.Info("user logged in", "user_id", result.Identity.ID, "method", method)
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "authenticated",
		"user":   result.Identity,
		"tokens": tokens,
	})
}

func (s *AuthServer) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	tokens, err := s.tokenManager.RefreshToken(req.RefreshToken)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "invalid or expired refresh token"})
		return
	}

	writeJSON(w, http.StatusOK, tokens)
}

func (s *AuthServer) HandleLogout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	s.tokenManager.RevokeSession(req.RefreshToken)
	writeJSON(w, http.StatusOK, map[string]any{"status": "logged_out"})
}

func (s *AuthServer) HandleMe(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": "missing token"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := s.tokenManager.ValidateToken(token)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id": claims.Sub,
		"email":   claims.Email,
		"phone":   claims.Phone,
		"name":    claims.Name,
		"role":    claims.Role,
	})
}

func (s *AuthServer) HandleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	writeJSON(w, http.StatusOK, map[string]any{"status": "reset_link_sent", "email": req.Email})
}

func (s *AuthServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "auth-server"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

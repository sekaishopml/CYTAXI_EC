package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/sekaishopml/cytaxi/backend/mobile"
)

type MobileServer struct {
	manager *mobile.Manager
	logger  *slog.Logger
}

func NewMobileServer(manager *mobile.Manager, logger *slog.Logger) *MobileServer {
	return &MobileServer{manager: manager, logger: logger}
}

func (s *MobileServer) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"user_id"`
		DeviceType string `json:"device_type"`
		Token      string `json:"token"`
		AppVersion string `json:"app_version"`
		OSVersion  string `json:"os_version"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	device := s.manager.RegisterDevice(req.DeviceType, req.UserID, req.Token, req.AppVersion, req.OSVersion)
	s.logger.Info("device registered", "id", device.ID, "type", req.DeviceType)
	writeJSON(w, http.StatusCreated, device)
}

func (s *MobileServer) HandleDevices(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	devices := s.manager.GetDevices(userID)
	writeJSON(w, http.StatusOK, map[string]any{"devices": devices, "count": len(devices)})
}

func (s *MobileServer) HandleSync(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var action mobile.OfflineAction
		json.NewDecoder(r.Body).Decode(&action)
		result := s.manager.EnqueueAction(action)
		writeJSON(w, http.StatusAccepted, result)
		return
	}

	userID := r.URL.Query().Get("user_id")
	results := s.manager.SyncActions(userID)
	if results == nil {
		writeJSON(w, http.StatusOK, map[string]any{"synced": []mobile.SyncResult{}, "message": "No pending actions"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"synced": results, "count": len(results)})
}

func (s *MobileServer) HandlePush(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID string            `json:"user_id"`
		Title  string            `json:"title"`
		Body   string            `json:"body"`
		Data   map[string]any    `json:"data"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	pn := s.manager.SendPush(req.UserID, req.Title, req.Body, req.Data)
	s.logger.Info("push notification sent", "id", pn.ID, "user", req.UserID)
	writeJSON(w, http.StatusCreated, pn)
}

func (s *MobileServer) HandleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			UserID   string `json:"user_id"`
			DeviceID string `json:"device_id"`
			TTLHours int    `json:"ttl_hours"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.TTLHours == 0 { req.TTLHours = 24 }

		session := s.manager.CreateSession(req.UserID, req.DeviceID, time.Duration(req.TTLHours)*time.Hour)
		writeJSON(w, http.StatusCreated, session)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	session, err := s.manager.ValidateSession(sessionID)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error(), "valid": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"session": session, "valid": true})
}

func (s *MobileServer) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"device_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	s.manager.RevokeDevice(req.DeviceID)
	writeJSON(w, http.StatusOK, map[string]any{"status": "revoked"})
}

func (s *MobileServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	status := s.manager.GetStatus()
	status["status"] = "ok"
	status["service"] = "mobile-platform"
	writeJSON(w, http.StatusOK, status)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

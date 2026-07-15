package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/infrastructure/experience"
)

type DxServer struct {
	manager *experience.Manager
	logger  *slog.Logger
}

func NewDxServer(manager *experience.Manager, logger *slog.Logger) *DxServer {
	return &DxServer{manager: manager, logger: logger}
}

func (s *DxServer) HandleEarnings(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	if r.Method == "POST" {
		var req struct {
			Amount float64 `json:"amount"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		result := s.manager.UpdateEarnings(driverID, req.Amount)
		writeJSON(w, http.StatusOK, result)
		return
	}

	earnings := s.manager.GetEarnings(driverID)
	writeJSON(w, http.StatusOK, earnings)
}

func (s *DxServer) HandleBonuses(w http.ResponseWriter, r *http.Request) {
	bonuses := s.manager.GetBonuses()
	writeJSON(w, http.StatusOK, map[string]any{"bonuses": bonuses, "count": len(bonuses)})
}

func (s *DxServer) HandleGoals(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	if r.Method == "POST" {
		var g experience.DriverGoal
		json.NewDecoder(r.Body).Decode(&g)
		result := s.manager.SetGoal(driverID, g)
		writeJSON(w, http.StatusCreated, result)
		return
	}

	goals := s.manager.GetGoals(driverID)
	writeJSON(w, http.StatusOK, map[string]any{"goals": goals, "count": len(goals)})
}

func (s *DxServer) HandleShifts(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	if r.Method == "POST" {
		var req struct {
			Action string `json:"action"` // start, end
		}
		json.NewDecoder(r.Body).Decode(&req)

		if req.Action == "start" {
			shift := s.manager.StartShift(driverID)
			s.logger.Info("shift started", "driver", driverID, "shift", shift.ID)
			writeJSON(w, http.StatusCreated, shift)
		} else {
			shift := s.manager.EndShift(driverID)
			if shift == nil {
				writeJSON(w, http.StatusBadRequest, map[string]any{"error": "no active shift"})
				return
			}
			s.logger.Info("shift ended", "driver", driverID)
			writeJSON(w, http.StatusOK, shift)
		}
		return
	}

	shifts := s.manager.GetShifts(driverID)
	writeJSON(w, http.StatusOK, map[string]any{"shifts": shifts, "count": len(shifts)})
}

func (s *DxServer) HandlePerformance(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	if r.Method == "POST" {
		var req struct {
			AcceptRate   float64 `json:"acceptance_rate"`
			CompleteRate float64 `json:"completion_rate"`
			CancelRate   float64 `json:"cancel_rate"`
			Rating       float64 `json:"rating"`
			Trips        int     `json:"trip_count"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		s.manager.UpdatePerformance(driverID, req.AcceptRate, req.CompleteRate, req.CancelRate, req.Rating, req.Trips)
		writeJSON(w, http.StatusOK, map[string]any{"status": "updated"})
		return
	}

	perf := s.manager.GetPerformance(driverID)
	writeJSON(w, http.StatusOK, perf)
}

func (s *DxServer) HandlePreferences(w http.ResponseWriter, r *http.Request) {
	driverID := r.PathValue("driver_id")

	if r.Method == "POST" {
		var p experience.DriverPreference
		json.NewDecoder(r.Body).Decode(&p)
		result := s.manager.UpdatePreferences(driverID, p)
		writeJSON(w, http.StatusOK, result)
		return
	}

	prefs := s.manager.GetPreferences(driverID)
	writeJSON(w, http.StatusOK, prefs)
}

func (s *DxServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "driver-experience"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

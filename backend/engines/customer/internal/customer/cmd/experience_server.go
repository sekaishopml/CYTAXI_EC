package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/customer/internal/customer/infrastructure/experience"
)

type CxServer struct {
	manager *experience.Manager
	logger  *slog.Logger
}

func NewCxServer(manager *experience.Manager, logger *slog.Logger) *CxServer {
	return &CxServer{manager: manager, logger: logger}
}

func (s *CxServer) HandleFavorites(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")

	if r.Method == "POST" {
		var fav experience.FavoritePlace
		json.NewDecoder(r.Body).Decode(&fav)
		result := s.manager.AddFavorite(customerID, fav)
		writeJSON(w, http.StatusCreated, result)
		return
	}

	favs := s.manager.GetFavorites(customerID)
	writeJSON(w, http.StatusOK, map[string]any{"favorites": favs, "count": len(favs)})
}

func (s *CxServer) HandleSavedPlaces(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")

	if r.Method == "POST" {
		var route experience.SavedRoute
		json.NewDecoder(r.Body).Decode(&route)
		result := s.manager.SaveRoute(customerID, route)
		writeJSON(w, http.StatusCreated, result)
		return
	}

	routes := s.manager.GetSavedRoutes(customerID)
	writeJSON(w, http.StatusOK, map[string]any{"routes": routes, "count": len(routes)})
}

func (s *CxServer) HandlePreferences(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")

	if r.Method == "POST" {
		var prefs experience.CustomerPreference
		json.NewDecoder(r.Body).Decode(&prefs)
		result := s.manager.UpdatePreferences(customerID, prefs)
		writeJSON(w, http.StatusOK, result)
		return
	}

	prefs := s.manager.GetPreferences(customerID)
	writeJSON(w, http.StatusOK, prefs)
}

func (s *CxServer) HandleLoyalty(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")

	acc := s.manager.GetLoyalty(customerID)
	rewards := experience.GetRewards()
	writeJSON(w, http.StatusOK, map[string]any{
		"account": acc,
		"available_rewards": rewards,
	})
}

func (s *CxServer) HandleEarnPoints(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CustomerID string  `json:"customer_id"`
		Amount     float64 `json:"amount"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	acc := s.manager.EarnPoints(req.CustomerID, req.Amount)
	s.logger.Info("points earned", "customer", req.CustomerID, "points", acc.Points)
	writeJSON(w, http.StatusOK, acc)
}

func (s *CxServer) HandleSupport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			CustomerID  string `json:"customer_id"`
			Subject     string `json:"subject"`
			Description string `json:"description"`
			Category    string `json:"category"`
			Priority    string `json:"priority"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		ticket := s.manager.CreateTicket(req.CustomerID, req.Subject, req.Description, req.Category, req.Priority)
		s.logger.Info("support ticket created", "id", ticket.ID)
		writeJSON(w, http.StatusCreated, ticket)
		return
	}

	customerID := r.URL.Query().Get("customer_id")
	tickets := s.manager.GetTickets(customerID)
	writeJSON(w, http.StatusOK, map[string]any{"tickets": tickets, "count": len(tickets)})
}

func (s *CxServer) HandleNotifications(w http.ResponseWriter, r *http.Request) {
	customerID := r.PathValue("customer_id")

	if r.Method == "POST" {
		var prefs experience.NotificationPref
		json.NewDecoder(r.Body).Decode(&prefs)
		result := s.manager.UpdateNotificationPrefs(customerID, prefs)
		writeJSON(w, http.StatusOK, result)
		return
	}

	prefs := s.manager.GetNotificationPrefs(customerID)
	writeJSON(w, http.StatusOK, prefs)
}

func (s *CxServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "customer-experience"})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

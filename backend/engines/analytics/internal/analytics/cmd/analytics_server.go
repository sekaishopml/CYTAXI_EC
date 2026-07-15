package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/infrastructure/analytics"
)

type AnalyticsServer struct {
	manager *analytics.Manager
	logger  *slog.Logger
}

func NewAnalyticsServer(manager *analytics.Manager, logger *slog.Logger) *AnalyticsServer {
	return &AnalyticsServer{manager: manager, logger: logger}
}

func (s *AnalyticsServer) HandleDashboard(w http.ResponseWriter, r *http.Request) {
	dashType := r.URL.Query().Get("type")
	if dashType == "" { dashType = "general" }

	dashboard := s.manager.GetDashboard(dashType)

	// Simulate events for live dashboard
	consumer := s.manager.GetConsumer()
	consumer.ConsumeTripCompleted(12.50)
	consumer.ConsumeTripCompleted(8.00)
	consumer.ConsumeRating(5)

	writeJSON(w, http.StatusOK, dashboard)
}

func (s *AnalyticsServer) HandleKPIs(w http.ResponseWriter, r *http.Request) {
	snap := s.manager.GetSnapshot()
	writeJSON(w, http.StatusOK, map[string]any{
		"kpis": []map[string]any{
			{"name": "Total Trips", "value": snap.TotalTrips, "trend": "up"},
			{"name": "Completed Trips", "value": snap.CompletedTrips},
			{"name": "Cancel Rate", "value": snap.CancelRate, "unit": "%"},
			{"name": "Total Revenue", "value": snap.TotalRevenue, "unit": "USD"},
			{"name": "Active Drivers", "value": snap.ActiveDrivers},
			{"name": "Registered Users", "value": snap.RegisteredUsers},
			{"name": "Avg Rating", "value": snap.AvgRating, "unit": "/5"},
			{"name": "Peak Hour", "value": snap.PeakHour, "unit": "hour"},
			{"name": "Promos Used", "value": snap.PromosUsed},
			{"name": "Coupons Used", "value": snap.CouponsUsed},
		},
	})
}

func (s *AnalyticsServer) HandleReports(w http.ResponseWriter, r *http.Request) {
	reportType := r.URL.Query().Get("type")
	if reportType == "" { reportType = "daily" }

	period := r.URL.Query().Get("period")
	if period == "" { period = "today" }

	report := s.manager.GetReport(reportType, period)
	writeJSON(w, http.StatusOK, report)
}

func (s *AnalyticsServer) HandleTrends(w http.ResponseWriter, r *http.Request) {
	snap := s.manager.GetSnapshot()
	consumer := s.manager.GetConsumer()

	writeJSON(w, http.StatusOK, map[string]any{
		"metric": "trips",
		"trends": map[string]any{
			"trips_today":     snap.TotalTrips,
			"revenue_today":   snap.TotalRevenue,
			"active_drivers":  snap.ActiveDrivers,
			"by_hour":         consumer.ByHour(),
		},
	})
}

func (s *AnalyticsServer) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	snap := s.manager.GetSnapshot()
	writeJSON(w, http.StatusOK, snap)
}

func (s *AnalyticsServer) HandleExport(w http.ResponseWriter, r *http.Request) {
	format := r.URL.Query().Get("format")
	if format == "" { format = "json" }

	snap := s.manager.GetSnapshot()
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=cytaxi_report.json")
	json.NewEncoder(w).Encode(snap)
}

func (s *AnalyticsServer) HandleEvent(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Type   string  `json:"type"`
		Amount float64 `json:"amount,omitempty"`
		Score  int     `json:"score,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	consumer := s.manager.GetConsumer()
	switch req.Type {
	case "TripCompleted":
		consumer.ConsumeTripCompleted(req.Amount)
	case "TripCancelled":
		consumer.ConsumeTripCancelled()
	case "PaymentCompleted":
		consumer.ConsumePaymentCompleted(req.Amount)
	case "Refund":
		consumer.ConsumeRefund(req.Amount)
	case "DriverOnline":
		consumer.ConsumeDriverOnline()
	case "UserRegistered":
		consumer.ConsumeUserRegistered()
	case "PromoUsed":
		consumer.ConsumePromoUsed()
	case "CouponUsed":
		consumer.ConsumeCouponUsed()
	case "Rating":
		consumer.ConsumeRating(req.Score)
	}

	s.logger.Info("event consumed", "type", req.Type)
	writeJSON(w, http.StatusOK, map[string]any{"status": "consumed"})
}

func (s *AnalyticsServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	snap := s.manager.GetSnapshot()
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok", "service": "analytics-engine",
		"events_processed": snap.TotalTrips,
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

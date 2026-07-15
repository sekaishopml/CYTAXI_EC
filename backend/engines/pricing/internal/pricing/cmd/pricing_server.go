package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/pricing/internal/pricing/infrastructure/pricing"
)

type PricingServer struct {
	registry *pricing.StrategyRegistry
	rules    *pricing.RuleManager
	logger   *slog.Logger
	history  sync.Map // trip_id → FareResult
}

func NewPricingServer(registry *pricing.StrategyRegistry, rules *pricing.RuleManager, logger *slog.Logger) *PricingServer {
	return &PricingServer{registry: registry, rules: rules, logger: logger}
}

func (s *PricingServer) HandleCalculate(w http.ResponseWriter, r *http.Request) {
	var input pricing.FareInput
	json.NewDecoder(r.Body).Decode(&input)

	strategyName := s.rules.GetActiveStrategy()
	strategy, err := s.registry.Get(strategyName)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	result, err := strategy.Calculate(input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	// Apply promo/coupon
	if input.PromoCode != "" {
		s.rules.ApplyPromotion(result, input.PromoCode)
	}
	if input.CouponCode != "" {
		s.rules.ApplyCoupon(result, input.CouponCode)
	}

	tripID := fmt.Sprintf("trip_%d", time.Now().UnixNano())
	s.history.Store(tripID, result)

	s.logger.Info("fare calculated", "strategy", strategyName, "total", result.Total)

	writeJSON(w, http.StatusOK, map[string]any{
		"trip_id":  tripID,
		"strategy": strategyName,
		"fare":     result,
	})
}

func (s *PricingServer) HandleEstimate(w http.ResponseWriter, r *http.Request) {
	var input pricing.FareInput
	json.NewDecoder(r.Body).Decode(&input)

	strategy, _ := s.registry.Get("distance")
	result, err := strategy.Calculate(input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"estimate": result,
		"strategy": "distance",
		"note":     "Estimated fare - final may vary",
	})
}

func (s *PricingServer) HandleGetRules(w http.ResponseWriter, r *http.Request) {
	active := s.rules.GetActiveStrategy()
	writeJSON(w, http.StatusOK, map[string]any{
		"active_strategy": active,
		"available":       []string{"fixed", "distance", "dynamic", "zone"},
	})
}

func (s *PricingServer) HandlePromotions(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var promo pricing.Promotion
		json.NewDecoder(r.Body).Decode(&promo)
		s.rules.AddPromotion(&promo)
		writeJSON(w, http.StatusCreated, promo)
		return
	}

	s.logger.Info("promotions list")
	writeJSON(w, http.StatusOK, map[string]any{"message": "Use POST to create promotions"})
}

func (s *PricingServer) HandleCoupon(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CouponCode string  `json:"coupon_code"`
		Fare       float64 `json:"fare_subtotal"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	result := &pricing.FareResult{Subtotal: req.Fare, Total: req.Fare}
	err := s.rules.ApplyCoupon(result, req.CouponCode)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error(), "valid": false})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"valid": true, "discount": result.CouponDiscount, "final": result.Total})
}

func (s *PricingServer) HandleHistory(w http.ResponseWriter, r *http.Request) {
	var history []any
	s.history.Range(func(k, v any) bool {
		history = append(history, map[string]any{"trip_id": k, "fare": v})
		return true
	})
	writeJSON(w, http.StatusOK, map[string]any{"history": history, "count": len(history)})
}

func (s *PricingServer) HandleSetStrategy(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Strategy string `json:"strategy"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	if _, err := s.registry.Get(req.Strategy); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}
	s.rules.SetActiveStrategy(req.Strategy)
	s.logger.Info("pricing strategy changed", "strategy", req.Strategy)
	writeJSON(w, http.StatusOK, map[string]any{"status": "updated", "strategy": req.Strategy})
}

func (s *PricingServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status": "ok", "service": "dynamic-pricing",
		"active_strategy": s.rules.GetActiveStrategy(),
	})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

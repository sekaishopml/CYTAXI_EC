package metrics

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type BusinessMetrics struct {
	ID           valueobject.MetricID       `json:"id"`
	DateRange    valueobject.DateRange      `json:"date_range"`
	Trips        TripMetrics               `json:"trip_metrics"`
	Revenue      RevenueMetrics            `json:"revenue_metrics"`
	Drivers      DriverMetrics             `json:"driver_metrics"`
	Customers    CustomerMetrics            `json:"customer_metrics"`
	Platform     PlatformMetrics           `json:"platform_metrics"`
	CalculatedAt time.Time                 `json:"calculated_at"`
}

type TripMetrics struct {
	TotalTrips    int     `json:"total_trips"`
	Completed     int     `json:"completed"`
	Cancelled     int     `json:"cancelled"`
	CancelRate    float64 `json:"cancel_rate_pct"`
	AvgWaitTime   float64 `json:"avg_wait_seconds"`
	AvgTripTime   float64 `json:"avg_trip_seconds"`
	AvgDistanceKM float64 `json:"avg_distance_km"`
}

type RevenueMetrics struct {
	TotalRevenue float64 `json:"total_revenue"`
	AvgFare      float64 `json:"avg_fare"`
	Commission   float64 `json:"commission"`
	Payouts      float64 `json:"driver_payouts"`
	PlatformFee  float64 `json:"platform_fee"`
}

type DriverMetrics struct {
	ActiveDrivers   int     `json:"active_drivers"`
	IdleDrivers     int     `json:"idle_drivers"`
	AcceptanceRate  float64 `json:"acceptance_rate_pct"`
	AvgRating       float64 `json:"avg_rating"`
	TopDrivers      int     `json:"top_drivers"`
}

type CustomerMetrics struct {
	ActiveUsers    int     `json:"active_users"`
	NewUsers       int     `json:"new_users"`
	FrequentUsers  int     `json:"frequent_users"`
	AvgTripsPerUser float64 `json:"avg_trips_per_user"`
}

type PlatformMetrics struct {
	Uptime          float64 `json:"uptime_pct"`
	MatchingTime    float64 `json:"avg_matching_seconds"`
	SuccessRate     float64 `json:"success_rate_pct"`
	PromoCouponsUsed int    `json:"promo_coupons_used"`
}

func NewBusinessMetrics(dr valueobject.DateRange) *BusinessMetrics {
	return &BusinessMetrics{
		ID:           valueobject.NewMetricID(),
		DateRange:    dr,
		CalculatedAt: time.Now(),
	}
}

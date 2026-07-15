package pricing

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type FareInput struct {
	DistanceKM    float64 `json:"distance_km"`
	DurationSec   int     `json:"duration_sec"`
	IsNight       bool    `json:"is_night"`
	IsAirport     bool    `json:"is_airport"`
	DemandLevel   int     `json:"demand_level"` // 1-10
	SupplyLevel   int     `json:"supply_level"` // 1-10
	Zone          string  `json:"zone"`
	VehicleType   string  `json:"vehicle_type"`
	IsCorporate   bool    `json:"is_corporate"`
	HasPromo      bool    `json:"has_promo"`
	PromoCode     string  `json:"promo_code,omitempty"`
	CouponCode    string  `json:"coupon_code,omitempty"`
	WaitingSec    int     `json:"waiting_sec"`
}

type FareResult struct {
	BaseFare      float64 `json:"base_fare"`
	DistanceFare  float64 `json:"distance_fare"`
	TimeFare      float64 `json:"time_fare"`
	NightCharge   float64 `json:"night_charge,omitempty"`
	AirportCharge float64 `json:"airport_charge,omitempty"`
	DemandFactor  float64 `json:"demand_factor"`
	SupplyFactor  float64 `json:"supply_factor"`
	ZoneCharge    float64 `json:"zone_charge,omitempty"`
	VehicleExtra  float64 `json:"vehicle_extra,omitempty"`
	WaitingFee    float64 `json:"waiting_fee,omitempty"`
	Subtotal      float64 `json:"subtotal"`
	PromoDiscount float64 `json:"promo_discount,omitempty"`
	CouponDiscount float64 `json:"coupon_discount,omitempty"`
	Tax           float64 `json:"tax"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
	Strategy      string  `json:"strategy"`
}

type PricingStrategy interface {
	Name() string
	Calculate(input FareInput) (*FareResult, error)
}

type FixedFareStrategy struct{}

func (s *FixedFareStrategy) Name() string { return "fixed" }
func (s *FixedFareStrategy) Calculate(input FareInput) (*FareResult, error) {
	result := &FareResult{
		BaseFare: 3.00, Currency: "USD", Strategy: "fixed",
		Subtotal: 3.00, Total: 3.00 + 3.00*0.12,
		Tax: 3.00 * 0.12,
	}
	return result, nil
}

type DistanceFareStrategy struct {
	RatePerKM float64
	MinFare   float64
}

func (s *DistanceFareStrategy) Name() string { return "distance" }
func (s *DistanceFareStrategy) Calculate(input FareInput) (*FareResult, error) {
	base := 1.00
	dist := math.Max(input.DistanceKM, 0.1) * s.RatePerKM
	time := float64(input.DurationSec) * 0.03
	subtotal := base + dist + time

	if subtotal < s.MinFare { subtotal = s.MinFare + 0.50 }

	result := &FareResult{
		BaseFare: base, DistanceFare: dist, TimeFare: time,
		Subtotal: subtotal, Tax: subtotal * 0.12,
		Total: subtotal + subtotal*0.12, Currency: "USD", Strategy: "distance",
	}
	return result, nil
}

type DynamicFareStrategy struct {
	BaseRate       float64
	DemandMultiplier float64
	SupplyDivider  float64
}

func (s *DynamicFareStrategy) Name() string { return "dynamic" }
func (s *DynamicFareStrategy) Calculate(input FareInput) (*FareResult, error) {
	base := s.BaseRate
	dist := input.DistanceKM * 0.60
	time := float64(input.DurationSec) * 0.03

	demandFactor := 1.0 + float64(input.DemandLevel-5)*s.DemandMultiplier
	if demandFactor < 0.8 { demandFactor = 0.8 }
	if demandFactor > 3.0 { demandFactor = 3.0 }

	supplyFactor := 1.0
	if input.SupplyLevel > 0 {
		supplyFactor = 1.0 + (10.0-float64(input.SupplyLevel))*s.SupplyDivider
	}

	subtotal := (base + dist + time) * demandFactor * supplyFactor

	if input.IsNight {
		subtotal += 2.50
	}
	if input.IsAirport {
		subtotal += 3.00
	}

	waiting := float64(input.WaitingSec) * 0.02

	result := &FareResult{
		BaseFare: base, DistanceFare: dist, TimeFare: time,
		DemandFactor: demandFactor, SupplyFactor: supplyFactor,
		NightCharge: boolToFloat(input.IsNight, 2.50),
		AirportCharge: boolToFloat(input.IsAirport, 3.00),
		WaitingFee: waiting,
		Subtotal: subtotal + waiting, Tax: (subtotal + waiting) * 0.12,
		Total: (subtotal + waiting) * 1.12,
		Currency: "USD", Strategy: "dynamic",
	}
	return result, nil
}

type ZoneFareStrategy struct {
	ZoneRates map[string]float64
}

func (s *ZoneFareStrategy) Name() string { return "zone" }
func (s *ZoneFareStrategy) Calculate(input FareInput) (*FareResult, error) {
	zoneRate, ok := s.ZoneRates[input.Zone]
	if !ok { zoneRate = 1.0 }

	base := 1.50
	dist := input.DistanceKM * 0.50 * zoneRate
	time := float64(input.DurationSec) * 0.025
	subtotal := base + dist + time

	result := &FareResult{
		BaseFare: base, DistanceFare: dist, TimeFare: time,
		ZoneCharge: (zoneRate - 1.0) * base * 5,
		Subtotal: subtotal, Tax: subtotal * 0.12,
		Total: subtotal * 1.12, Currency: "USD", Strategy: "zone",
	}
	return result, nil
}

type StrategyRegistry struct {
	strategies map[string]PricingStrategy
}

func NewRegistry() *StrategyRegistry {
	r := &StrategyRegistry{strategies: make(map[string]PricingStrategy)}
	r.Register(&FixedFareStrategy{})
	r.Register(&DistanceFareStrategy{RatePerKM: 0.50, MinFare: 3.00})
	r.Register(&DynamicFareStrategy{BaseRate: 1.00, DemandMultiplier: 0.15, SupplyDivider: 0.1})
	r.Register(&ZoneFareStrategy{ZoneRates: map[string]float64{"downtown": 1.0, "airport": 1.3, "suburb": 0.9}})
	return r
}

func (r *StrategyRegistry) Register(s PricingStrategy) { r.strategies[s.Name()] = s }
func (r *StrategyRegistry) Get(name string) (PricingStrategy, error) {
	s, ok := r.strategies[name]
	if !ok { return nil, fmt.Errorf("strategy %s not found", name) }
	return s, nil
}

type Promotion struct {
	Code      string    `json:"code"`
	Type      string    `json:"type"` // percentage, fixed, free_ride
	Value     float64   `json:"value"`
	MinFare   float64   `json:"min_fare"`
	MaxDiscount float64  `json:"max_discount"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Active    bool      `json:"active"`
}

type Coupon struct {
	Code      string    `json:"code"`
	Type      string    `json:"type"`
	Value     float64   `json:"value"`
	MaxUses   int       `json:"max_uses"`
	UsedCount int       `json:"used_count"`
	MinFare   float64   `json:"min_fare"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Active    bool      `json:"active"`
}

type RuleManager struct {
	mu          sync.RWMutex
	promotions   map[string]*Promotion
	coupons      map[string]*Coupon
	activeStrategy string
}

func NewRuleManager() *RuleManager {
	return &RuleManager{
		promotions:    make(map[string]*Promotion),
		coupons:       make(map[string]*Coupon),
		activeStrategy: "dynamic",
	}
}

func (m *RuleManager) AddPromotion(p *Promotion) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.promotions[p.Code] = p
}

func (m *RuleManager) GetPromotion(code string) *Promotion {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.promotions[code]
}

func (m *RuleManager) AddCoupon(c *Coupon) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.coupons[c.Code] = c
}

func (m *RuleManager) GetCoupon(code string) *Coupon {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.coupons[code]
}

func (m *RuleManager) SetActiveStrategy(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeStrategy = name
}

func (m *RuleManager) GetActiveStrategy() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.activeStrategy
}

func (m *RuleManager) ApplyPromotion(fare *FareResult, promoCode string) error {
	p := m.GetPromotion(promoCode)
	if p == nil || !p.Active { return fmt.Errorf("promotion not found or inactive") }
	if time.Now().Before(p.ValidFrom) || time.Now().After(p.ValidTo) { return fmt.Errorf("promotion expired") }

	switch p.Type {
	case "percentage":
		discount := fare.Subtotal * p.Value / 100
		if p.MaxDiscount > 0 && discount > p.MaxDiscount { discount = p.MaxDiscount }
		fare.PromoDiscount = discount
		fare.Total -= discount
	case "fixed":
		fare.PromoDiscount = p.Value
		fare.Total -= p.Value
	case "free_ride":
		fare.PromoDiscount = fare.Total
		fare.Total = 0
	}
	if fare.Total < 0 { fare.Total = 0 }
	return nil
}

func (m *RuleManager) ApplyCoupon(fare *FareResult, couponCode string) error {
	c := m.GetCoupon(couponCode)
	if c == nil || !c.Active { return fmt.Errorf("coupon not found or inactive") }
	if c.MaxUses > 0 && c.UsedCount >= c.MaxUses { return fmt.Errorf("coupon exhausted") }
	if c.MinFare > 0 && fare.Subtotal < c.MinFare { return fmt.Errorf("minimum fare not met") }

	c.UsedCount++

	switch c.Type {
	case "percentage":
		fare.CouponDiscount = fare.Subtotal * c.Value / 100
	case "fixed":
		fare.CouponDiscount = c.Value
	}
	fare.Total -= fare.CouponDiscount
	if fare.Total < 0 { fare.Total = 0 }
	return nil
}

func boolToFloat(b bool, val float64) float64 { if b { return val }; return 0 }

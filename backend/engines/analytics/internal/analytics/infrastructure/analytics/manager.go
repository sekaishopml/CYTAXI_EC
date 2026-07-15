package analytics

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type KPI struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Current  float64 `json:"current"`
	Previous float64 `json:"previous"`
	Change   float64 `json:"change_pct"`
	Target   float64 `json:"target"`
	Status   string  `json:"status"` // ok, warning, critical
	UpdatedAt time.Time `json:"updated_at"`
}

type Dashboard struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"` // operations, commercial, quality, infrastructure
	Title      string          `json:"title"`
	KPIs       []KPI           `json:"kpis"`
	Charts     []Chart         `json:"charts"`
	TopDrivers []TopItem       `json:"top_drivers,omitempty"`
	TopZones   []TopItem       `json:"top_zones,omitempty"`
	GeneratedAt time.Time      `json:"generated_at"`
}

type Chart struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Type   string        `json:"type"` // line, bar, pie, gauge, area
	Labels []string      `json:"labels"`
	Series []ChartSeries `json:"series"`
}

type ChartSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

type TopItem struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Rank  int     `json:"rank"`
}

type Report struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Type        string    `json:"type"` // daily, weekly, monthly, custom
	Period      string    `json:"period"`
	Summary     string    `json:"summary"`
	Metrics     []Metric  `json:"metrics"`
	GeneratedAt time.Time `json:"generated_at"`
}

type Metric struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Unit      string  `json:"unit"`
	Trend     string  `json:"trend"`     // up, down, stable
	ChangePct float64 `json:"change_pct"`
}

type Trend struct {
	Metric    string      `json:"metric"`
	Period    string      `json:"period"`
	DataPoints []DataPoint `json:"data_points"`
}

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

type EventConsumer struct {
	mu        sync.RWMutex
	trips     int64
	completed int64
	cancelled int64
	revenue   float64
	payments  int64
	refunds   float64
	driversOnline int64
	driversOffline int64
	usersRegistered int64
	promosUsed int64
	couponsUsed int64
	ratings   []float64
	byHour    [24]int64
	byDay     [7]int64
}

func NewEventConsumer() *EventConsumer { return &EventConsumer{} }

func (c *EventConsumer) ConsumeTripCompleted(fare float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.trips++
	c.completed++
	c.revenue += fare
	c.byHour[time.Now().Hour()]++
	c.byDay[time.Now().Weekday()]++
}

func (c *EventConsumer) ConsumeTripCancelled() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.trips++
	c.cancelled++
}

func (c *EventConsumer) ConsumePaymentCompleted(amount float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.payments++
	c.revenue += amount
}

func (c *EventConsumer) ConsumeRefund(amount float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.refunds += amount
}

func (c *EventConsumer) ConsumeDriverOnline() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.driversOnline++
}

func (c *EventConsumer) ConsumeUserRegistered() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.usersRegistered++
}

func (c *EventConsumer) ConsumePromoUsed() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.promosUsed++
}

func (c *EventConsumer) ConsumeCouponUsed() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.couponsUsed++
}

func (c *EventConsumer) ConsumeRating(score int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ratings = append(c.ratings, float64(score))
	if len(c.ratings) > 1000 {
		c.ratings = c.ratings[1:]
	}
}

func (c *EventConsumer) GetSnapshot() AnalyticsSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()

	avgRating := 0.0
	for _, r := range c.ratings {
		avgRating += r
	}
	if len(c.ratings) > 0 { avgRating /= float64(len(c.ratings)) }

	cancelRate := 0.0
	if c.trips > 0 { cancelRate = float64(c.cancelled) / float64(c.trips) * 100 }

	peakHour := 0
	for i := 1; i < 24; i++ {
		if c.byHour[i] > c.byHour[peakHour] { peakHour = i }
	}

	return AnalyticsSnapshot{
		TotalTrips:       c.trips,
		CompletedTrips:   c.completed,
		CancelledTrips:   c.cancelled,
		CancelRate:       math.Round(cancelRate*100) / 100,
		TotalRevenue:     math.Round(c.revenue*100) / 100,
		TotalPayments:    c.payments,
		TotalRefunds:     math.Round(c.refunds*100) / 100,
		ActiveDrivers:    c.driversOnline,
		RegisteredUsers:  c.usersRegistered,
		PromosUsed:       c.promosUsed,
		CouponsUsed:      c.couponsUsed,
		AvgRating:        math.Round(avgRating*100) / 100,
		PeakHour:         peakHour,
		SnapshotAt:       time.Now(),
	}
}

type AnalyticsSnapshot struct {
	TotalTrips      int64     `json:"total_trips"`
	CompletedTrips  int64     `json:"completed_trips"`
	CancelledTrips  int64     `json:"cancelled_trips"`
	CancelRate      float64   `json:"cancel_rate_pct"`
	TotalRevenue    float64   `json:"total_revenue"`
	TotalPayments   int64     `json:"total_payments"`
	TotalRefunds    float64   `json:"total_refunds"`
	ActiveDrivers   int64     `json:"active_drivers"`
	RegisteredUsers int64     `json:"registered_users"`
	PromosUsed      int64     `json:"promos_used"`
	CouponsUsed     int64     `json:"coupons_used"`
	AvgRating       float64   `json:"avg_rating"`
	PeakHour        int       `json:"peak_hour"`
	SnapshotAt      time.Time `json:"snapshot_at"`
}

type Manager struct {
	consumer *EventConsumer
}

func NewManager() *Manager {
	return &Manager{consumer: NewEventConsumer()}
}

func (m *Manager) GetConsumer() *EventConsumer { return m.consumer }

func (m *Manager) GetDashboard(dashboardType string) *Dashboard {
	snap := m.consumer.GetSnapshot()

	d := &Dashboard{
		ID: fmt.Sprintf("dash_%d", time.Now().Unix()), Type: dashboardType,
		Title: "Executive Dashboard", GeneratedAt: time.Now(),
	}

	switch dashboardType {
	case "operations":
		d.KPIs = []KPI{
			{Name: "Total Trips", Category: "ops", Current: float64(snap.TotalTrips), Previous: float64(snap.TotalTrips-10), Change: 15.0, Target: float64(snap.TotalTrips)+10, Status: "ok"},
			{Name: "Cancel Rate", Category: "ops", Current: snap.CancelRate, Previous: 12.0, Change: -20.0, Target: 5.0, Status: "warning"},
			{Name: "Active Drivers", Category: "ops", Current: float64(snap.ActiveDrivers), Target: 50, Status: "ok"},
		}
	case "commercial":
		d.KPIs = []KPI{
			{Name: "Revenue", Category: "commercial", Current: snap.TotalRevenue, Previous: snap.TotalRevenue*0.9, Change: 10.0, Target: snap.TotalRevenue*1.2, Status: "ok"},
			{Name: "Avg Rating", Category: "quality", Current: snap.AvgRating, Target: 4.5, Status: "ok"},
			{Name: "Promos Used", Category: "commercial", Current: float64(snap.PromosUsed), Target: 100, Status: "ok"},
		}
	default:
		d.KPIs = []KPI{
			{Name: "Trips", Category: "ops", Current: float64(snap.TotalTrips), Status: "ok"},
			{Name: "Revenue", Category: "commercial", Current: snap.TotalRevenue, Status: "ok"},
			{Name: "Drivers", Category: "ops", Current: float64(snap.ActiveDrivers), Status: "ok"},
			{Name: "Users", Category: "growth", Current: float64(snap.RegisteredUsers), Status: "ok"},
		}
	}

	d.Charts = []Chart{
		{ID: "trips_by_hour", Title: "Trips by Hour", Type: "bar",
			Labels: []string{"00","01","02","03","04","05","06","07","08","09","10","11","12","13","14","15","16","17","18","19","20","21","22","23"},
			Series: []ChartSeries{{Name: "Trips", Data: int64ToFloat64(m.consumer.byHour[:])}},
		},
	}

	return d
}

func (m *Manager) GetReport(reportType, period string) *Report {
	snap := m.consumer.GetSnapshot()
	return &Report{
		ID: fmt.Sprintf("rpt_%d", time.Now().Unix()), Title: reportType + " Report",
		Type: reportType, Period: period, GeneratedAt: time.Now(),
		Summary: fmt.Sprintf("Report covering %d trips with $%.2f revenue", snap.TotalTrips, snap.TotalRevenue),
		Metrics: []Metric{
			{Name: "Trips", Value: float64(snap.TotalTrips), Unit: "count", Trend: "up"},
			{Name: "Revenue", Value: snap.TotalRevenue, Unit: "USD", Trend: "up"},
			{Name: "Active Drivers", Value: float64(snap.ActiveDrivers), Unit: "count", Trend: "stable"},
			{Name: "Cancel Rate", Value: snap.CancelRate, Unit: "%", Trend: "down"},
		},
	}
}

func (m *Manager) GetSnapshot() AnalyticsSnapshot { return m.consumer.GetSnapshot() }

func int64ToFloat64(arr []int64) []float64 {
	result := make([]float64, len(arr))
	for i, v := range arr { result[i] = float64(v) }
	return result
}

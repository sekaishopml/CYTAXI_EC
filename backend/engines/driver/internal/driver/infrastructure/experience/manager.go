package experience

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type EarningsSummary struct {
	DriverID    string    `json:"driver_id"`
	Today       float64   `json:"today"`
	ThisWeek    float64   `json:"this_week"`
	ThisMonth   float64   `json:"this_month"`
	Total       float64   `json:"total"`
	TripsToday  int       `json:"trips_today"`
	TripsTotal  int       `json:"trips_total"`
	HoursOnline float64   `json:"hours_online_today"`
	AvgPerHour  float64   `json:"avg_per_hour"`
	AvgPerTrip  float64   `json:"avg_per_trip"`
	LastUpdated time.Time `json:"last_updated"`
}

type BonusProgram struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Target      int     `json:"target"` // trips count
	Reward      float64 `json:"reward"`
	Active      bool    `json:"active"`
}

type DriverGoal struct {
	ID        string    `json:"id"`
	DriverID  string    `json:"driver_id"`
	Type      string    `json:"type"` // trips, revenue, hours, rating
	Target    float64   `json:"target"`
	Current   float64   `json:"current"`
	Progress  float64   `json:"progress_pct"`
	Deadline  time.Time `json:"deadline"`
	Completed bool      `json:"completed"`
}

type ShiftRecord struct {
	ID        string     `json:"id"`
	DriverID  string     `json:"driver_id"`
	StartAt   time.Time  `json:"start_at"`
	EndAt     *time.Time  `json:"end_at,omitempty"`
	TripsDone int        `json:"trips_done"`
	Earnings  float64    `json:"earnings"`
	Status    string     `json:"status"` // active, completed, cancelled
}

type PerformanceMetric struct {
	DriverID       string  `json:"driver_id"`
	AcceptanceRate float64 `json:"acceptance_rate"`
	CompletionRate float64 `json:"completion_rate"`
	CancelRate     float64 `json:"cancel_rate"`
	AvgRating      float64 `json:"avg_rating"`
	ResponseTime   float64 `json:"avg_response_sec"`
	TripCount      int     `json:"trip_count"`
	Rank           int     `json:"rank,omitempty"`
	Trend          string  `json:"trend"` // improving, stable, declining
}

type DriverPreference struct {
	DriverID        string  `json:"driver_id"`
	MaxDistance     float64 `json:"max_distance_km"`
	MinFare         float64 `json:"min_fare"`
	PreferredZones  []string `json:"preferred_zones"`
	AutoAccept      bool    `json:"auto_accept"`
	QuietMode       bool    `json:"quiet_mode"`
	NavigationApp   string  `json:"navigation_app"`
	Language        string  `json:"language"`
}

type Manager struct {
	mu         sync.RWMutex
	earnings   map[string]*EarningsSummary
	goals      map[string][]DriverGoal
	shifts     map[string][]ShiftRecord
	bonuses    []BonusProgram
	perf       map[string]*PerformanceMetric
	prefs      map[string]*DriverPreference
}

func NewManager() *Manager {
	return &Manager{
		earnings: make(map[string]*EarningsSummary),
		goals:    make(map[string][]DriverGoal),
		shifts:   make(map[string][]ShiftRecord),
		perf:     make(map[string]*PerformanceMetric),
		prefs:    make(map[string]*DriverPreference),
		bonuses:  defaultBonuses(),
	}
}

func defaultBonuses() []BonusProgram {
	return []BonusProgram{
		{ID: "bon_1", Name: "First 10 Trips", Description: "Complete your first 10 trips", Target: 10, Reward: 25, Active: true},
		{ID: "bon_2", Name: "Weekend Warrior", Description: "Complete 20 trips on weekend", Target: 20, Reward: 50, Active: true},
		{ID: "bon_3", Name: "Night Owl", Description: "Complete 15 trips between 10PM-4AM", Target: 15, Reward: 35, Active: true},
		{ID: "bon_4", Name: "100 Trip Club", Description: "100 lifetime trips", Target: 100, Reward: 100, Active: true},
	}
}

func (m *Manager) GetEarnings(driverID string) *EarningsSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	e, ok := m.earnings[driverID]
	if !ok {
		return &EarningsSummary{DriverID: driverID, LastUpdated: time.Now()}
	}
	return e
}

func (m *Manager) UpdateEarnings(driverID string, amount float64) *EarningsSummary {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.earnings[driverID]
	if !ok {
		e = &EarningsSummary{DriverID: driverID}
		m.earnings[driverID] = e
	}
	e.Today += amount
	e.ThisWeek += amount
	e.ThisMonth += amount
	e.Total += amount
	e.TripsToday++
	e.TripsTotal++
	if e.HoursOnline > 0 { e.AvgPerHour = e.Today / e.HoursOnline }
	e.LastUpdated = time.Now()
	return e
}

func (m *Manager) GetBonuses() []BonusProgram {
	return m.bonuses
}

func (m *Manager) SetGoal(driverID string, g DriverGoal) *DriverGoal {
	g.ID = fmt.Sprintf("gl_%d", time.Now().UnixNano())
	g.DriverID = driverID
	m.mu.Lock()
	m.goals[driverID] = append(m.goals[driverID], g)
	m.mu.Unlock()
	return &g
}

func (m *Manager) GetGoals(driverID string) []DriverGoal {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.goals[driverID]
}

func (m *Manager) StartShift(driverID string) *ShiftRecord {
	s := &ShiftRecord{
		ID: fmt.Sprintf("sft_%d", time.Now().UnixNano()),
		DriverID: driverID, StartAt: time.Now(), Status: "active",
	}
	m.mu.Lock()
	m.shifts[driverID] = append(m.shifts[driverID], *s)
	m.mu.Unlock()
	return s
}

func (m *Manager) EndShift(driverID string) *ShiftRecord {
	m.mu.Lock()
	defer m.mu.Unlock()
	shifts := m.shifts[driverID]
	for i := len(shifts) - 1; i >= 0; i-- {
		if shifts[i].Status == "active" {
			now := time.Now()
			shifts[i].EndAt = &now
			shifts[i].Status = "completed"
			m.shifts[driverID] = shifts
			return &shifts[i]
		}
	}
	return nil
}

func (m *Manager) GetShifts(driverID string) []ShiftRecord {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.shifts[driverID]
}

func (m *Manager) GetPerformance(driverID string) *PerformanceMetric {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.perf[driverID]
	if !ok {
		return &PerformanceMetric{DriverID: driverID, Trend: "stable"}
	}
	return p
}

func (m *Manager) UpdatePerformance(driverID string, acceptRate, completeRate, cancelRate, rating float64, trips int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	p := &PerformanceMetric{
		DriverID: driverID, AcceptanceRate: acceptRate, CompletionRate: completeRate,
		CancelRate: cancelRate, AvgRating: rating, TripCount: trips,
		Trend: "stable",
	}
	m.perf[driverID] = p
}

func (m *Manager) UpdatePreferences(driverID string, p DriverPreference) *DriverPreference {
	p.DriverID = driverID
	m.mu.Lock()
	m.prefs[driverID] = &p
	m.mu.Unlock()
	return &p
}

func (m *Manager) GetPreferences(driverID string) *DriverPreference {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.prefs[driverID]
	if !ok {
		return &DriverPreference{
			DriverID: driverID, MaxDistance: 15, MinFare: 3.00,
			PreferredZones: []string{"downtown"}, Language: "es",
		}
	}
	return p
}

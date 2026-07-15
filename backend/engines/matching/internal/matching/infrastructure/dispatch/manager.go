package dispatch

import (
	"fmt"
	"math"
	"sort"
	"sync"
	"time"
)

type DispatchRequest struct {
	ID         string    `json:"id"`
	TripID     string    `json:"trip_id"`
	PickupLat  float64   `json:"pickup_lat"`
	PickupLng  float64   `json:"pickup_lng"`
	Zone       string    `json:"zone"`
	Priority   int       `json:"priority"` // 1-10
	MaxRadius  float64   `json:"max_radius_meters"`
	MaxRetries int       `json:"max_retries"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Candidate struct {
	DriverID    string  `json:"driver_id"`
	Name        string  `json:"name"`
	Lat         float64 `json:"lat"`
	Lng         float64 `json:"lng"`
	Distance    float64 `json:"distance_meters"`
	ETA         int     `json:"eta_seconds"`
	Rating      float64 `json:"rating"`
	AcceptRate  float64 `json:"acceptance_rate"`
	Zone        string  `json:"zone"`
	Score       float64 `json:"score"`
	Rank        int     `json:"rank"`
}

type DispatchZone struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	RadiusKM float64 `json:"radius_km"`
	Active   bool    `json:"active"`
	Priority int     `json:"priority"`
}

type QueuePolicy string

const (
	PolicyFIFO     QueuePolicy = "fifo"
	PolicyPriority QueuePolicy = "priority"
	PolicyBalanced QueuePolicy = "balanced"
	PolicyNearest   QueuePolicy = "nearest_driver"
)

type ScoringConfig struct {
	DistanceWeight    float64 `json:"distance_weight"`
	ETAWeight         float64 `json:"eta_weight"`
	RatingWeight      float64 `json:"rating_weight"`
	AcceptRateWeight  float64 `json:"accept_rate_weight"`
	ZoneWeight        float64 `json:"zone_weight"`
}

func DefaultScoring() ScoringConfig {
	return ScoringConfig{
		DistanceWeight: 0.35, ETAWeight: 0.30, RatingWeight: 0.15,
		AcceptRateWeight: 0.10, ZoneWeight: 0.10,
	}
}

type AssignmentAttempt struct {
	ID         string    `json:"id"`
	RequestID  string    `json:"request_id"`
	DriverID   string    `json:"driver_id"`
	Attempt    int       `json:"attempt"`
	Result     string    `json:"result"` // offered, accepted, rejected, expired, failed
	CreatedAt  time.Time `json:"created_at"`
}

type Metrics struct {
	mu              sync.RWMutex
	TotalRequests   int64 `json:"total_requests"`
	TotalAssignments int64 `json:"total_assignments"`
	TotalRejections  int64 `json:"total_rejections"`
	TotalRetries     int64 `json:"total_retries"`
	TotalCancellations int64 `json:"total_cancellations"`
	AvgAssignmentMs  float64 `json:"avg_assignment_ms"`
	AcceptanceRate   float64 `json:"acceptance_rate_pct"`
	AvgAttempts      float64 `json:"avg_attempts"`
	AvgETASeconds    float64 `json:"avg_eta_seconds"`
	activeRequests   map[string]time.Time
}

func NewMetrics() *Metrics {
	return &Metrics{activeRequests: make(map[string]time.Time)}
}

func (m *Metrics) RecordRequest(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRequests++
	m.activeRequests[id] = time.Now()
}

func (m *Metrics) RecordAssignment(requestID string, eta int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalAssignments++
	if start, ok := m.activeRequests[requestID]; ok {
		elapsed := time.Since(start).Milliseconds()
		m.AvgAssignmentMs = (m.AvgAssignmentMs*float64(m.TotalAssignments-1) + float64(elapsed)) / float64(m.TotalAssignments)
		delete(m.activeRequests, requestID)
	}
	m.AvgETASeconds = (m.AvgETASeconds*float64(m.TotalAssignments-1) + float64(eta)) / float64(m.TotalAssignments)
	m.AcceptanceRate = float64(m.TotalAssignments) / float64(m.TotalRequests) * 100
}

func (m *Metrics) RecordRetry() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRetries++
}

func (m *Metrics) RecordRejection() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalRejections++
}

func (m *Metrics) RecordCancellation() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TotalCancellations++
}

type Manager struct {
	mu         sync.RWMutex
	requests   sync.Map
	zones      sync.Map
	attempts   sync.Map
	scoring    ScoringConfig
	policy     QueuePolicy
	metrics    *Metrics
}

func NewManager() *Manager {
	return &Manager{
		scoring: DefaultScoring(),
		policy:  PolicyBalanced,
		metrics: NewMetrics(),
	}
}

func (m *Manager) SetScoring(cfg ScoringConfig) { m.scoring = cfg }
func (m *Manager) SetPolicy(p QueuePolicy)      { m.policy = p }
func (m *Manager) GetMetrics() *Metrics         { return m.metrics }

func (m *Manager) AddZone(zone DispatchZone) { m.zones.Store(zone.ID, zone) }

func (m *Manager) StartDispatch(req DispatchRequest) *DispatchRequest {
	req.ID = fmt.Sprintf("dsp_%d", time.Now().UnixNano())
	req.Status = "pending"
	req.CreatedAt = time.Now()
	if req.MaxRadius == 0 { req.MaxRadius = 5000 }
	if req.MaxRetries == 0 { req.MaxRetries = 3 }
	m.requests.Store(req.ID, &req)
	m.metrics.RecordRequest(req.ID)
	return &req
}

func (m *Manager) GetRequest(id string) (*DispatchRequest, error) {
	v, ok := m.requests.Load(id)
	if !ok { return nil, fmt.Errorf("dispatch %s not found", id) }
	return v.(*DispatchRequest), nil
}

func (m *Manager) ScoreCandidates(request *DispatchRequest, drivers []Candidate) []Candidate {
	for i := range drivers {
		distScore := 1.0
		if request.MaxRadius > 0 && drivers[i].Distance > 0 {
			distScore = 1.0 - (drivers[i].Distance / request.MaxRadius)
		}
		if distScore < 0 { distScore = 0 }

		etaScore := 1.0
		if drivers[i].ETA > 0 {
			etaScore = 1.0 - math.Min(float64(drivers[i].ETA)/600.0, 1.0) // max 10min
		}

		zoneBonus := 1.0
		if drivers[i].Zone == request.Zone { zoneBonus = 2.0 }

		drivers[i].Score = 
			m.scoring.DistanceWeight*distScore +
			m.scoring.ETAWeight*etaScore +
			m.scoring.RatingWeight*(drivers[i].Rating/5.0) +
			m.scoring.AcceptRateWeight*(drivers[i].AcceptRate/100.0) +
			m.scoring.ZoneWeight*zoneBonus
	}

	sort.Slice(drivers, func(i, j int) bool {
		return drivers[i].Score > drivers[j].Score
	})

	for i := range drivers {
		drivers[i].Rank = i + 1
	}
	return drivers
}

func (m *Manager) GetBestCandidate(drivers []Candidate) *Candidate {
	if len(drivers) == 0 { return nil }
	return &drivers[0]
}

func (m *Manager) RecordAttempt(requestID, driverID string, attempt int, result string) *AssignmentAttempt {
	a := &AssignmentAttempt{
		ID: fmt.Sprintf("att_%d", time.Now().UnixNano()),
		RequestID: requestID, DriverID: driverID,
		Attempt: attempt, Result: result, CreatedAt: time.Now(),
	}
	m.attempts.Store(a.ID, a)

	switch result {
	case "accepted":
		m.metrics.RecordAssignment(requestID, 0)
	case "rejected":
		m.metrics.RecordRejection()
	case "expired":
		m.metrics.RecordRetry()
	}
	return a
}

func (m *Manager) RetryDispatch(id string) (*DispatchRequest, error) {
	req, err := m.GetRequest(id)
	if err != nil { return nil, err }

	attempts := m.countAttempts(id)
	if attempts >= req.MaxRetries {
		req.Status = "failed"
		m.requests.Store(id, req)
		m.metrics.RecordCancellation()
		return req, fmt.Errorf("max retries (%d) exceeded", req.MaxRetries)
	}

	req.Status = "retrying"
	m.requests.Store(id, req)
	m.metrics.RecordRetry()
	return req, nil
}

func (m *Manager) CompleteDispatch(id, driverID string) error {
	req, err := m.GetRequest(id)
	if err != nil { return err }
	req.Status = "completed"
	m.requests.Store(id, req)
	return nil
}

func (m *Manager) CancelDispatch(id string) error {
	req, err := m.GetRequest(id)
	if err != nil { return err }
	req.Status = "cancelled"
	m.requests.Store(id, req)
	m.metrics.RecordCancellation()
	return nil
}

func (m *Manager) countAttempts(requestID string) int {
	count := 0
	m.attempts.Range(func(_, v any) bool {
		if v.(*AssignmentAttempt).RequestID == requestID {
			count++
		}
		return true
	})
	return count
}

func (m *Manager) GetCandidatesForZone(zoneID string, radius float64) []Candidate {
	// Mock candidates based on zone
	names := []string{"Carlos M.", "Ana P.", "Luis R.", "Jose V."}
	zones := []string{zoneID, zoneID, "zone_b", zoneID}

	var candidates []Candidate
	for i, name := range names {
		distance := 200 + float64(i)*350
		candidates = append(candidates, Candidate{
			DriverID: fmt.Sprintf("drv_%d", 1000+i),
			Name: name, Lat: -0.18 + float64(i)*0.001, Lng: -78.47 + float64(i)*0.001,
			Distance: distance, ETA: 60 + i*40, Rating: 4.0 + float64(i)*0.3,
			AcceptRate: 70 + float64(i)*8, Zone: zones[i], Rank: i + 1,
		})
	}
	return candidates
}

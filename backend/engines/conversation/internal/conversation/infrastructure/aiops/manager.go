package aiops

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type Incident struct {
	ID          string    `json:"id"`
	Type        string    `json:"type"` // service_down, high_latency, payment_failure, matching_failure, db_error
	Severity    string    `json:"severity"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
	Status      string    `json:"status"` // detected, analyzing, resolved
	DetectedAt  time.Time `json:"detected_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type Recommendation struct {
	ID          string    `json:"id"`
	IncidentID  string    `json:"incident_id,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // performance, security, reliability, cost
	Priority    string    `json:"priority"`
	Action      string    `json:"action"`
	Accepted    bool      `json:"accepted"`
	Executed    bool      `json:"executed"`
	GeneratedAt time.Time `json:"generated_at"`
}

type Runbook struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	IncidentType string  `json:"incident_type"`
	Severity    string   `json:"severity"`
	Steps       []string `json:"steps"`
	AutoFix     bool     `json:"auto_fix"`
}

type KnowledgeEntry struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
	Views   int      `json:"views"`
}

type Automation struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Trigger     string    `json:"trigger"` // health_degraded, error_spike, resource_low
	Action      string    `json:"action"`
	Enabled     bool      `json:"enabled"`
	ExecutedAt  time.Time `json:"executed_at"`
	Result      string    `json:"result"`
}

type Manager struct {
	mu              sync.RWMutex
	incidents       map[string]*Incident
	recommendations map[string][]*Recommendation
	runbooks        map[string]*Runbook
	knowledge       map[string]*KnowledgeEntry
	automations     map[string]*Automation
	acceptedRecs    int
	executedRecs    int
	totalRecs       int
}

func NewManager() *Manager {
	m := &Manager{
		incidents:       make(map[string]*Incident),
		recommendations: make(map[string][]*Recommendation),
		runbooks:        make(map[string]*Runbook),
		knowledge:       make(map[string]*KnowledgeEntry),
		automations:     make(map[string]*Automation),
	}
	m.initRunbooks()
	m.initKnowledge()
	m.initAutomations()
	return m
}

func (m *Manager) initRunbooks() {
	m.runbooks["service_down"] = &Runbook{
		ID: "rb_1", Title: "Service Down Response", IncidentType: "service_down", Severity: "critical",
		Steps: []string{
			"1. Check Gateway health: curl /health",
			"2. Check affected engine: curl {engine}:{port}/health",
			"3. Check Docker: docker ps | grep {service}",
			"4. Check logs: docker logs {service} --tail=50",
			"5. Restart: docker restart {service}",
			"6. Verify recovery: repeat step 1-2",
		},
	}
	m.runbooks["high_latency"] = &Runbook{
		ID: "rb_2", Title: "High Latency Response", IncidentType: "high_latency", Severity: "high",
		Steps: []string{
			"1. Check rate limiter: /health endpoint latency",
			"2. Check PostgreSQL: pg_isready -U cytaxi",
			"3. Check Redis: redis-cli PING",
			"4. Check CPU/Memory: docker stats {service}",
			"5. Scale if needed: increase replicas",
		},
	}
	m.runbooks["payment_failure"] = &Runbook{
		ID: "rb_3", Title: "Payment Failure Response", IncidentType: "payment_failure", Severity: "high",
		Steps: []string{
			"1. Verify Payment Engine: /payments/health",
			"2. Check payment provider status: /payments/providers",
			"3. Review recent payments: /payments/history",
			"4. Trigger retry if needed: POST /payments/retry",
		},
	}
	m.runbooks["matching_failure"] = &Runbook{
		ID: "rb_4", Title: "Matching Failure Response", IncidentType: "matching_failure", Severity: "medium",
		Steps: []string{
			"1. Verify Matching Engine: /matching/health",
			"2. Verify Driver Engine: /driver/health",
			"3. Check active drivers: GET /dispatch/metrics",
			"4. Restart matching if needed",
		},
	}
}

func (m *Manager) initKnowledge() {
	m.knowledge["kb_1"] = &KnowledgeEntry{
		ID: "kb_1", Title: "How to restart a service", Content: "docker restart {service_name}", Tags: []string{"docker", "restart"},
	}
	m.knowledge["kb_2"] = &KnowledgeEntry{
		ID: "kb_2", Title: "Database connection troubleshooting", Content: "Check pg_isready, verify credentials in .env, check network: docker network ls", Tags: []string{"database", "postgres"},
	}
	m.knowledge["kb_3"] = &KnowledgeEntry{
		ID: "kb_3", Title: "Scaling guidelines", Content: "API Gateway: 2+ replicas. Trip Engine: 2 replicas for 1000+ trips/hour. Add more as load increases.", Tags: []string{"scaling", "performance"},
	}
}

func (m *Manager) initAutomations() {
	m.automations["auto_1"] = &Automation{
		ID: "auto_1", Name: "Auto-restart unhealthy services", Description: "Restart any service that fails health check 3x in 5min", Trigger: "health_degraded", Enabled: false,
	}
	m.automations["auto_2"] = &Automation{
		ID: "auto_2", Name: "Auto-scale on high CPU", Description: "Increase replicas when CPU > 80% for 5min", Trigger: "resource_low", Enabled: false,
	}
}

func (m *Manager) DetectIncident(incidentType, severity, source, description string) *Incident {
	i := &Incident{
		ID: fmt.Sprintf("inc_%d", time.Now().UnixNano()),
		Type: incidentType, Severity: severity, Source: source,
		Description: description, Status: "detected", DetectedAt: time.Now(),
	}
	m.mu.Lock()
	m.incidents[i.ID] = i
	m.mu.Unlock()

	// Auto-generate recommendation
	rec := m.GenerateRecommendation(i.ID, "Auto-generated: "+description, "reliability", "high",
		"Review runbook for "+incidentType)
	rec.IncidentID = i.ID

	return i
}

func (m *Manager) GetRunbook(incidentType string) *Runbook {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.runbooks[incidentType]
}

func (m *Manager) GenerateRecommendation(incidentID, description, category, priority, action string) *Recommendation {
	r := &Recommendation{
		ID: fmt.Sprintf("rec_%d", time.Now().UnixNano()),
		IncidentID: incidentID, Title: "Recommendation",
		Description: description, Category: category,
		Priority: priority, Action: action, GeneratedAt: time.Now(),
	}
	m.mu.Lock()
	m.totalRecs++
	m.recommendations[incidentID] = append(m.recommendations[incidentID], r)
	m.mu.Unlock()
	return r
}

func (m *Manager) AcceptRecommendation(recID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, recs := range m.recommendations {
		for _, r := range recs {
			if r.ID == recID {
				r.Accepted = true
				m.acceptedRecs++
				return nil
			}
		}
	}
	return fmt.Errorf("recommendation not found")
}

func (m *Manager) GetIncidents() []Incident {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []Incident
	for _, i := range m.incidents {
		result = append(result, *i)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].DetectedAt.After(result[j].DetectedAt) })
	return result
}

func (m *Manager) GetRecommendations() []*Recommendation {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*Recommendation
	for _, recs := range m.recommendations {
		result = append(result, recs...)
	}
	return result
}

func (m *Manager) GetStatus() map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	openIncidents := 0
	for _, i := range m.incidents {
		if i.Status == "detected" || i.Status == "analyzing" {
			openIncidents++
		}
	}

	acceptRate := 0.0
	if m.totalRecs > 0 {
		acceptRate = float64(m.acceptedRecs) / float64(m.totalRecs) * 100
	}

	return map[string]any{
		"total_incidents":      len(m.incidents),
		"open_incidents":       openIncidents,
		"total_recommendations": m.totalRecs,
		"accepted":             m.acceptedRecs,
		"acceptance_rate_pct":  acceptRate,
		"runbooks_available":   len(m.runbooks),
		"automations_enabled":  m.countEnabledAutomations(),
	}
}

func (m *Manager) countEnabledAutomations() int {
	count := 0
	for _, a := range m.automations {
		if a.Enabled { count++ }
	}
	return count
}

func (m *Manager) GetKnowledge(query string) []KnowledgeEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []KnowledgeEntry
	for _, k := range m.knowledge {
		for _, tag := range k.Tags {
			if tag == query || query == "" {
				k.Views++
				result = append(result, *k)
				break
			}
		}
	}
	return result
}

package trust

import (
	"fmt"
	"math"
	"sync"
	"time"
)

type Rating struct {
	ID        string    `json:"id"`
	TripID    string    `json:"trip_id"`
	FromID    string    `json:"from_id"`
	ToID      string    `json:"to_id"`
	Score     int       `json:"score"` // 1-5
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Review struct {
	ID        string    `json:"id"`
	RatingID  string    `json:"rating_id"`
	Moderated bool      `json:"moderated"`
	Visible   bool      `json:"visible"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Incident struct {
	ID          string    `json:"id"`
	ReportedBy  string    `json:"reported_by"`
	ReportedID  string    `json:"reported_id"`
	TripID      string    `json:"trip_id,omitempty"`
	Type        IncidentType `json:"type"`
	Severity    string    `json:"severity"` // low, medium, high, critical
	Status      string    `json:"status"`   // open, investigating, resolved, closed
	Description string    `json:"description"`
	Resolution  string    `json:"resolution,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type IncidentType string

const (
	IncidentSafety    IncidentType = "safety"
	IncidentBehavior  IncidentType = "behavior"
	IncidentFraud     IncidentType = "fraud"
	IncidentVehicle   IncidentType = "vehicle"
	IncidentPayment   IncidentType = "payment"
	IncidentOther     IncidentType = "other"
)

type Appeal struct {
	ID          string    `json:"id"`
	IncidentID  string    `json:"incident_id"`
	FromID      string    `json:"from_id"`
	Reason      string    `json:"reason"`
	Status      string    `json:"status"` // pending, approved, rejected
	ReviewedBy  string    `json:"reviewed_by,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type TrustScore struct {
	UserID    string             `json:"user_id"`
	Score     float64            `json:"score"`  // 0-100
	Level     string             `json:"level"`  // excellent, good, fair, poor, critical
	Components TrustScoreComponents `json:"components"`
	History   []ScoreChange      `json:"history,omitempty"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type TrustScoreComponents struct {
	RatingsScore     float64 `json:"ratings_score"`
	CompletionRate   float64 `json:"completion_rate"`
	AcceptanceRate   float64 `json:"acceptance_rate,omitempty"`
	IncidentScore    float64 `json:"incident_score"`
	VerificationScore float64 `json:"verification_score"`
}

type ScoreChange struct {
	OldScore float64   `json:"old_score"`
	NewScore float64   `json:"new_score"`
	Reason   string    `json:"reason"`
	At       time.Time `json:"at"`
}

type Manager struct {
	ratings    sync.Map
	incidents  sync.Map
	appeals    sync.Map
	scores     sync.Map
	reviews     sync.Map
}

func NewManager() *Manager { return &Manager{} }

func (m *Manager) SubmitRating(tripID, fromID, toID string, score int, comment string) *Rating {
	r := &Rating{
		ID: fmt.Sprintf("rat_%d", time.Now().UnixNano()),
		TripID: tripID, FromID: fromID, ToID: toID,
		Score: score, Comment: comment, CreatedAt: time.Now(),
	}
	m.ratings.Store(r.ID, r)

	review := &Review{ID: fmt.Sprintf("rev_%d", time.Now().UnixNano()), RatingID: r.ID, Visible: true}
	m.reviews.Store(review.ID, review)

	m.RecalculateTrustScore(toID)
	return r
}

func (m *Manager) GetRatingsByUser(userID string) []Rating {
	var result []Rating
	m.ratings.Range(func(_, v any) bool {
		r := v.(Rating)
		if r.ToID == userID {
			result = append(result, r)
		}
		return true
	})
	return result
}

func (m *Manager) ReportIncident(reportedBy, reportedID, tripID string, incidentType IncidentType, severity, description string) *Incident {
	i := &Incident{
		ID: fmt.Sprintf("inc_%d", time.Now().UnixNano()),
		ReportedBy: reportedBy, ReportedID: reportedID, TripID: tripID,
		Type: incidentType, Severity: severity, Status: "open",
		Description: description, CreatedAt: time.Now(),
	}
	m.incidents.Store(i.ID, i)
	return i
}

func (m *Manager) ResolveIncident(incidentID, resolution string) (*Incident, error) {
	v, ok := m.incidents.Load(incidentID)
	if !ok { return nil, fmt.Errorf("incident not found") }
	i := v.(*Incident)
	now := time.Now()
	i.Status = "resolved"
	i.Resolution = resolution
	i.ResolvedAt = &now
	m.incidents.Store(incidentID, i)

	m.RecalculateTrustScore(i.ReportedID)
	return i, nil
}

func (m *Manager) GetIncidentsByUser(userID string) []Incident {
	var result []Incident
	m.incidents.Range(func(_, v any) bool {
		i := v.(Incident)
		if i.ReportedID == userID {
			result = append(result, i)
		}
		return true
	})
	return result
}

func (m *Manager) FileAppeal(incidentID, fromID, reason string) *Appeal {
	a := &Appeal{
		ID: fmt.Sprintf("apl_%d", time.Now().UnixNano()),
		IncidentID: incidentID, FromID: fromID,
		Reason: reason, Status: "pending", CreatedAt: time.Now(),
	}
	m.appeals.Store(a.ID, a)
	return a
}

func (m *Manager) ResolveAppeal(appealID, reviewer string, approved bool) (*Appeal, error) {
	v, ok := m.appeals.Load(appealID)
	if !ok { return nil, fmt.Errorf("appeal not found") }
	a := v.(*Appeal)
	now := time.Now()
	if approved { a.Status = "approved" } else { a.Status = "rejected" }
	a.ReviewedBy = reviewer
	a.ResolvedAt = &now
	m.appeals.Store(appealID, a)
	return a, nil
}

func (m *Manager) GetTrustScore(userID string) *TrustScore {
	v, ok := m.scores.Load(userID)
	if ok { return v.(*TrustScore) }

	ts := &TrustScore{
		UserID: userID, Score: 100, Level: "excellent",
		Components: TrustScoreComponents{
			RatingsScore: 100, CompletionRate: 100, AcceptanceRate: 100,
			IncidentScore: 100, VerificationScore: 100,
		},
		UpdatedAt: time.Now(),
	}
	m.scores.Store(userID, ts)
	return ts
}

func (m *Manager) RecalculateTrustScore(userID string) *TrustScore {
	ratings := m.GetRatingsByUser(userID)
	incidents := m.GetIncidentsByUser(userID)

	var avgRating float64
	for _, r := range ratings {
		avgRating += float64(r.Score)
	}
	if len(ratings) > 0 {
		avgRating = (avgRating / float64(len(ratings))) * 20 // scale to 0-100
	} else {
		avgRating = 100
	}

	incidentScore := 100.0
	for _, i := range incidents {
		switch i.Severity {
		case "low": incidentScore -= 5
		case "medium": incidentScore -= 15
		case "high": incidentScore -= 30
		case "critical": incidentScore -= 50
		}
	}
	incidentScore = math.Max(0, incidentScore)

	score := avgRating*0.4 + 100*0.2 + incidentScore*0.4

	old, _ := m.scores.Load(userID)
	ts := &TrustScore{
		UserID: userID, Score: math.Round(score*100) / 100,
		Level: scoreToLevel(score),
		Components: TrustScoreComponents{
			RatingsScore: avgRating, IncidentScore: incidentScore,
		},
		UpdatedAt: time.Now(),
	}

	if old != nil {
		ts.History = old.(*TrustScore).History
		ts.History = append(ts.History, ScoreChange{
			OldScore: old.(*TrustScore).Score, NewScore: ts.Score,
			Reason: "recalculated", At: time.Now(),
		})
	}

	m.scores.Store(userID, ts)
	return ts
}

func (m *Manager) GetOpenIncidents() []Incident {
	var result []Incident
	m.incidents.Range(func(_, v any) bool {
		i := v.(Incident)
		if i.Status == "open" || i.Status == "investigating" {
			result = append(result, i)
		}
		return true
	})
	return result
}

func scoreToLevel(score float64) string {
	switch {
	case score >= 85: return "excellent"
	case score >= 70: return "good"
	case score >= 50: return "fair"
	case score >= 30: return "poor"
	default: return "critical"
	}
}

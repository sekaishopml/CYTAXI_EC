package candidate

import "time"

type CandidateID string

type Candidate struct {
	ID            CandidateID      `json:"id"`
	DriverID      string           `json:"driver_id"`
	DriverName    string           `json:"driver_name,omitempty"`
	Location      Location         `json:"location"`
	Vehicle       Vehicle          `json:"vehicle"`
	DistanceMeters float64          `json:"distance_meters"`
	ETA            time.Duration    `json:"eta"`
	Score         float64          `json:"score"`
	Status        CandidateStatus  `json:"status"`
	Reasons       []string         `json:"reasons,omitempty"`
}

type CandidateStatus string

const (
	CandidateAvailable CandidateStatus = "available"
	CandidateFiltered  CandidateStatus = "filtered"
	CandidateSelected  CandidateStatus = "selected"
	CandidateRejected  CandidateStatus = "rejected"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Vehicle struct {
	ID    string `json:"id"`
	Plate string `json:"plate"`
	Model string `json:"model"`
	Color string `json:"color"`
	Type  string `json:"type"`
}

type CandidateSet struct {
	TripID       string       `json:"trip_id"`
	Candidates   []Candidate  `json:"candidates"`
	GeneratedAt  time.Time    `json:"generated_at"`
	TotalCount   int          `json:"total_count"`
}

func (cs *CandidateSet) Add(c Candidate) {
	cs.Candidates = append(cs.Candidates, c)
	cs.TotalCount++
}

func (cs *CandidateSet) FilterByStatus(status CandidateStatus) []Candidate {
	var filtered []Candidate
	for _, c := range cs.Candidates {
		if c.Status == status {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

package candidate

import (
	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type CandidateStatus string

const (
	CandidatePending  CandidateStatus = "pending"
	CandidateEval     CandidateStatus = "evaluated"
	CandidateRanked   CandidateStatus = "ranked"
	CandidateSelected CandidateStatus = "selected"
	CandidateRejected CandidateStatus = "rejected"
)

type DriverCandidate struct {
	DriverID   valueobject.DriverID    `json:"driver_id"`
	Status     CandidateStatus         `json:"status"`
	Snapshot   valueobject.DriverSnapshot `json:"snapshot"`
	Rank       valueobject.CandidateRank `json:"rank"`
	Score      valueobject.MatchingScore `json:"score"`
	RejectReason string                  `json:"reject_reason,omitempty"`
}

type CandidateSet struct {
	MatchingID valueobject.MatchingID `json:"matching_id"`
	Candidates []DriverCandidate      `json:"candidates"`
	Total      int                    `json:"total"`
	Evaluated  int                    `json:"evaluated"`
	Ranked     int                    `json:"ranked"`
}

func NewCandidateSet(matchingID valueobject.MatchingID) *CandidateSet {
	return &CandidateSet{
		MatchingID: matchingID,
	}
}

func (cs *CandidateSet) Add(c DriverCandidate) {
	cs.Candidates = append(cs.Candidates, c)
	cs.Total++
}

func (cs *CandidateSet) SelectByRank(rank valueobject.CandidateRank) *DriverCandidate {
	for i := range cs.Candidates {
		if cs.Candidates[i].Rank == rank {
			return &cs.Candidates[i]
		}
	}
	return nil
}

func (cs *CandidateSet) TopCandidates(n int) []DriverCandidate {
	if n > len(cs.Candidates) {
		n = len(cs.Candidates)
	}
	return cs.Candidates[:n]
}

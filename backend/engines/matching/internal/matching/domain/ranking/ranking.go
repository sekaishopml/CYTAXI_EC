package ranking

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/matching/internal/matching/domain/valueobject"
)

type CandidateRanking struct {
	ID          valueobject.MatchingID `json:"id"`
	MatchingID  valueobject.MatchingID `json:"matching_id"`
	Candidates  []RankedCandidate      `json:"candidates"`
	Strategy    string                 `json:"strategy"`
	BestScore   float64                `json:"best_score"`
	CreatedAt   time.Time              `json:"created_at"`
}

type RankedCandidate struct {
	DriverID valueobject.DriverID    `json:"driver_id"`
	Rank     valueobject.CandidateRank `json:"rank"`
	Score    valueobject.MatchingScore `json:"score"`
	Distance float64                `json:"distance_meters"`
	ETA      int                    `json:"eta_seconds"`
	Priority valueobject.Priority   `json:"priority"`
}

func NewCandidateRanking(matchingID valueobject.MatchingID, strategy string) *CandidateRanking {
	return &CandidateRanking{
		ID:         matchingID,
		MatchingID: matchingID,
		Strategy:   strategy,
		CreatedAt:  time.Now(),
	}
}

func (cr *CandidateRanking) AddRanked(c RankedCandidate) {
	cr.Candidates = append(cr.Candidates, c)
	if c.Score.Float() > cr.BestScore {
		cr.BestScore = c.Score.Float()
	}
}

func (cr *CandidateRanking) Top() *RankedCandidate {
	if len(cr.Candidates) == 0 {
		return nil
	}
	return &cr.Candidates[0]
}

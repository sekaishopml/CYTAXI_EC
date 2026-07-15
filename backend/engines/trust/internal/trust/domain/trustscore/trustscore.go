package trustscore

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type TrustProfile struct {
	IdentityID   valueobject.IdentityID `json:"identity_id"`
	Score        float64                `json:"score"`
	Level        valueobject.TrustLevel  `json:"level"`
	Components   TrustComponents        `json:"components"`
	History      []TrustChange          `json:"history,omitempty"`
	CalculatedAt time.Time              `json:"calculated_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type TrustComponents struct {
	VerificationScore float64 `json:"verification_score"`
	ActivityScore     float64 `json:"activity_score"`
	CommunityScore    float64 `json:"community_score"`
	ComplianceScore   float64 `json:"compliance_score"`
}

type TrustChange struct {
	OldScore float64   `json:"old_score"`
	NewScore float64   `json:"new_score"`
	Reason   string    `json:"reason"`
	At       time.Time `json:"at"`
}

func NewTrustProfile(identityID valueobject.IdentityID) *TrustProfile {
	now := time.Now()
	return &TrustProfile{
		IdentityID:   identityID,
		Score:        0,
		Level:        valueobject.TrustNone,
		CalculatedAt: now,
		UpdatedAt:    now,
	}
}

func (tp *TrustProfile) Calculate() {
	comp := tp.Components
	tp.Score = (comp.VerificationScore*0.4 + comp.ActivityScore*0.3 + comp.CommunityScore*0.2 + comp.ComplianceScore*0.1)
	tp.Level = scoreToLevel(tp.Score)
	tp.CalculatedAt = time.Now()
	tp.UpdatedAt = time.Now()
}

func scoreToLevel(score float64) valueobject.TrustLevel {
	switch {
	case score >= 80:
		return valueobject.TrustPremium
	case score >= 60:
		return valueobject.TrustVerified
	case score >= 30:
		return valueobject.TrustBasic
	default:
		return valueobject.TrustNone
	}
}

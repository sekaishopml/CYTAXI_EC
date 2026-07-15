package fraud

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type FraudAssessment struct {
	ID          valueobject.FraudCheckID `json:"id"`
	IdentityID  valueobject.IdentityID   `json:"identity_id"`
	RiskLevel   valueobject.RiskLevel    `json:"risk_level"`
	Score       float64                  `json:"score"`
	Flags       []FraudFlag              `json:"flags,omitempty"`
	CheckedAt   time.Time                `json:"checked_at"`
}

type FraudFlag struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Severity    float64 `json:"severity"`
}

func NewFraudAssessment(identityID valueobject.IdentityID) *FraudAssessment {
	return &FraudAssessment{
		ID:         valueobject.NewFraudCheckID(),
		IdentityID: identityID,
		RiskLevel:  valueobject.RiskLow,
		CheckedAt:  time.Now(),
	}
}

func (f *FraudAssessment) AddFlag(flag FraudFlag) {
	f.Flags = append(f.Flags, flag)
	f.Score += flag.Severity
	f.recalculateRisk()
}

func (f *FraudAssessment) recalculateRisk() {
	switch {
	case f.Score >= 70:
		f.RiskLevel = valueobject.RiskCritical
	case f.Score >= 50:
		f.RiskLevel = valueobject.RiskHigh
	case f.Score >= 30:
		f.RiskLevel = valueobject.RiskMedium
	default:
		f.RiskLevel = valueobject.RiskLow
	}
}

type RiskAssessment struct {
	ID         valueobject.RiskAssessmentID `json:"id"`
	IdentityID valueobject.IdentityID       `json:"identity_id"`
	Level      valueobject.RiskLevel        `json:"level"`
	Details    map[string]any               `json:"details,omitempty"`
	AssessedAt time.Time                    `json:"assessed_at"`
}

func NewRiskAssessment(identityID valueobject.IdentityID) *RiskAssessment {
	return &RiskAssessment{
		ID:         valueobject.NewRiskAssessmentID(),
		IdentityID: identityID,
		Level:      valueobject.RiskLow,
		AssessedAt: time.Now(),
	}
}

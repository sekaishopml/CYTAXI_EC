package verification

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type Verification struct {
	ID         valueobject.VerificationID    `json:"id"`
	IdentityID valueobject.IdentityID        `json:"identity_id"`
	Type       VerificationType              `json:"type"`
	Status     valueobject.VerificationStatus `json:"status"`
	Provider   string                        `json:"provider,omitempty"`
	Result     VerificationResult            `json:"result,omitempty"`
	CreatedAt  time.Time                     `json:"created_at"`
	UpdatedAt  time.Time                     `json:"updated_at"`
}

type VerificationType string

const (
	VerifTypeDocument  VerificationType = "document"
	VerifTypeSelfie    VerificationType = "selfie"
	VerifTypeBiometric VerificationType = "biometric"
	VerifTypeKYC       VerificationType = "kyc"
	VerifTypeAML       VerificationType = "aml"
	VerifTypeAddress   VerificationType = "address"
	VerifTypePhone     VerificationType = "phone"
)

type VerificationResult struct {
	Passed   bool              `json:"passed"`
	Score    float64           `json:"score"`
	Details  map[string]string `json:"details,omitempty"`
}

func NewVerification(identityID valueobject.IdentityID, verifType VerificationType) *Verification {
	now := time.Now()
	return &Verification{
		ID:         valueobject.NewVerificationID(),
		IdentityID: identityID,
		Type:       verifType,
		Status:     valueobject.VerifPending,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (v *Verification) StartReview() {
	v.Status = valueobject.VerifInReview
	v.UpdatedAt = time.Now()
}

func (v *Verification) Approve(result VerificationResult) {
	v.Status = valueobject.VerifApproved
	v.Result = result
	v.UpdatedAt = time.Now()
}

func (v *Verification) Reject(reason string) {
	v.Status = valueobject.VerifRejected
	v.Result = VerificationResult{Passed: false, Details: map[string]string{"reason": reason}}
	v.UpdatedAt = time.Now()
}

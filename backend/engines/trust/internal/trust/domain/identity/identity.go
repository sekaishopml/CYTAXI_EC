package identity

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type Identity struct {
	ID            valueobject.IdentityID    `json:"id"`
	OwnerID       string                    `json:"owner_id"`
	Type          valueobject.IdentityType  `json:"type"`
	Status        valueobject.VerificationStatus `json:"status"`
	TrustLevel    valueobject.TrustLevel    `json:"trust_level"`
	Email         string                    `json:"email,omitempty"`
	Phone         string                    `json:"phone"`
	DocumentID    string                    `json:"document_id,omitempty"`
	VerifiedAt    *time.Time                `json:"verified_at,omitempty"`
	CreatedAt     time.Time                 `json:"created_at"`
	UpdatedAt     time.Time                 `json:"updated_at"`
}

func NewIdentity(ownerID string, idType valueobject.IdentityType, phone string) *Identity {
	now := time.Now()
	return &Identity{
		ID:         valueobject.NewIdentityID(),
		OwnerID:    ownerID,
		Type:       idType,
		Status:     valueobject.VerifPending,
		TrustLevel: valueobject.TrustNone,
		Phone:      phone,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func (i *Identity) Verify() {
	now := time.Now()
	i.Status = valueobject.VerifApproved
	i.TrustLevel = valueobject.TrustBasic
	i.VerifiedAt = &now
	i.UpdatedAt = now
}

func (i *Identity) Reject() {
	i.Status = valueobject.VerifRejected
	i.UpdatedAt = time.Now()
}

func (i *Identity) UpgradeTrust(level valueobject.TrustLevel) {
	i.TrustLevel = level
	i.UpdatedAt = time.Now()
}

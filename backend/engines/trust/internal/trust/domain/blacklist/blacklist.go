package blacklist

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type BlacklistEntry struct {
	IdentityID valueobject.IdentityID `json:"identity_id"`
	Reason     string                 `json:"reason"`
	Severity   string                 `json:"severity"`
	CreatedAt  time.Time              `json:"created_at"`
	ExpiresAt  *time.Time             `json:"expires_at,omitempty"`
}

func NewBlacklistEntry(identityID valueobject.IdentityID, reason, severity string) *BlacklistEntry {
	return &BlacklistEntry{
		IdentityID: identityID,
		Reason:     reason,
		Severity:   severity,
		CreatedAt:  time.Now(),
	}
}

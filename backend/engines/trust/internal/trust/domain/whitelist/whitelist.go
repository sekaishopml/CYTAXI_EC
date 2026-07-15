package whitelist

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type WhitelistEntry struct {
	IdentityID valueobject.IdentityID `json:"identity_id"`
	Reason     string                 `json:"reason"`
	CreatedAt  time.Time              `json:"created_at"`
}

func NewWhitelistEntry(identityID valueobject.IdentityID, reason string) *WhitelistEntry {
	return &WhitelistEntry{
		IdentityID: identityID,
		Reason:     reason,
		CreatedAt:  time.Now(),
	}
}

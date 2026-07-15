package document

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trust/internal/trust/domain/valueobject"
)

type Document struct {
	ID         valueobject.DocumentID     `json:"id"`
	IdentityID valueobject.IdentityID     `json:"identity_id"`
	Type       valueobject.DocumentType   `json:"type"`
	URL        string                     `json:"url"`
	Status     valueobject.VerificationStatus `json:"status"`
	ExpiresAt  *time.Time                 `json:"expires_at,omitempty"`
	UploadedAt time.Time                  `json:"uploaded_at"`
	VerifiedAt *time.Time                 `json:"verified_at,omitempty"`
}

func NewDocument(identityID valueobject.IdentityID, docType valueobject.DocumentType, url string) *Document {
	return &Document{
		ID:         valueobject.NewDocumentID(),
		IdentityID: identityID,
		Type:       docType,
		URL:        url,
		Status:     valueobject.VerifPending,
		UploadedAt: time.Now(),
	}
}

func (d *Document) Verify() {
	now := time.Now()
	d.Status = valueobject.VerifApproved
	d.VerifiedAt = &now
}

func (d *Document) Reject() {
	d.Status = valueobject.VerifRejected
}

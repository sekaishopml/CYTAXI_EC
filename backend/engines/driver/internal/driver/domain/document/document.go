package document

import (
	"fmt"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type DocumentType string

const (
	TypeDriverLicense DocumentType = "driver_license"
	TypeVehicleReg    DocumentType = "vehicle_registration"
	TypeInsurance     DocumentType = "insurance"
	TypeIdentityCard  DocumentType = "identity_card"
	TypeBackground    DocumentType = "background_check"
	TypeVehiclePhoto  DocumentType = "vehicle_photo"
)

type Document struct {
	ID         string               `json:"id"`
	DriverID   valueobject.DriverID `json:"driver_id"`
	Type       DocumentType         `json:"type"`
	URL        string               `json:"url"`
	Verified   bool                 `json:"verified"`
	ExpiresAt  *time.Time           `json:"expires_at,omitempty"`
	UploadedAt time.Time            `json:"uploaded_at"`
}

func NewDocument(driverID valueobject.DriverID, docType DocumentType, url string) *Document {
	return &Document{
		ID:         fmt.Sprintf("doc_%d", time.Now().UnixNano()),
		DriverID:   driverID,
		Type:       docType,
		URL:        url,
		UploadedAt: time.Now(),
	}
}

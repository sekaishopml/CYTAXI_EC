package license

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/driver/internal/driver/domain/valueobject"
)

type License struct {
	ID        valueobject.LicenseNumber `json:"id"`
	DriverID  valueobject.DriverID      `json:"driver_id"`
	Number    valueobject.LicenseNumber `json:"number"`
	Category  string                    `json:"category"`
	IssuedAt  time.Time                 `json:"issued_at"`
	ExpiresAt time.Time                 `json:"expires_at"`
	Verified  bool                      `json:"verified"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

func NewLicense(driverID valueobject.DriverID, number valueobject.LicenseNumber, category string) *License {
	now := time.Now()
	return &License{
		ID:        number,
		DriverID:  driverID,
		Number:    number,
		Category:  category,
		IssuedAt:  now,
		ExpiresAt: now.AddDate(5, 0, 0),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (l *License) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

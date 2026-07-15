package featureflag

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type FeatureFlag struct {
	ID          valueobject.FeatureFlagID `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Enabled     bool                      `json:"enabled"`
	Scope       string                    `json:"scope"` // global, city, percentage
	Percentage  int                       `json:"percentage,omitempty"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
}

func NewFeatureFlag(name, description string) *FeatureFlag {
	now := time.Now()
	return &FeatureFlag{
		ID:          valueobject.NewFeatureFlagID(),
		Name:        name,
		Description: description,
		Scope:       "global",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (f *FeatureFlag) Enable()  { f.Enabled = true; f.UpdatedAt = time.Now() }
func (f *FeatureFlag) Disable() { f.Enabled = false; f.UpdatedAt = time.Now() }

type SystemSetting struct {
	Key       valueobject.ConfigurationKey `json:"key"`
	Value     string                       `json:"value"`
	Group     string                       `json:"group"`
	UpdatedAt time.Time                    `json:"updated_at"`
}

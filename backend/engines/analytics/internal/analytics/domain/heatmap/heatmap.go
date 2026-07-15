package heatmap

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type Heatmap struct {
	ID          valueobject.HeatmapID `json:"id"`
	City        string                `json:"city"`
	DateRange   valueobject.DateRange `json:"date_range"`
	Granularity valueobject.TimeGranularity `json:"granularity"`
	Zones       []HeatZone            `json:"zones"`
	GeneratedAt time.Time             `json:"generated_at"`
}

type HeatZone struct {
	Lat       float64           `json:"lat"`
	Lng       float64           `json:"lng"`
	Level     valueobject.HeatLevel `json:"level"`
	Count     int               `json:"count"`
	Category  string            `json:"category"` // pickup, dropoff, both
}

func NewHeatmap(city string, dr valueobject.DateRange) *Heatmap {
	return &Heatmap{
		ID:          valueobject.NewHeatmapID(),
		City:        city,
		DateRange:   dr,
		Granularity: valueobject.GranHour,
		GeneratedAt: time.Now(),
	}
}

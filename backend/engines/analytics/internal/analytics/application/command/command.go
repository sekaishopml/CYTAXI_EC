package command

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type GenerateReport struct {
	Title       string
	Type        string
	DateRange   valueobject.DateRange
	Granularity valueobject.TimeGranularity
}

type CalculateMetrics struct {
	DateRange   valueobject.DateRange
	Granularity valueobject.TimeGranularity
}

type GenerateHeatmap struct {
	City      string
	DateRange valueobject.DateRange
}

type UpdateDashboard struct {
	KPIs   []any
	Charts []any
}

type CreateForecast struct {
	Metric      string
	Method      string
	Window      int
	Granularity valueobject.TimeGranularity
}

type RefreshAnalytics struct {
	Since time.Time
}

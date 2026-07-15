package query

import (
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type GetDashboard struct{}

type GetMetrics struct {
	DateRange   valueobject.DateRange
	Granularity valueobject.TimeGranularity
}

type GetRevenue struct {
	DateRange valueobject.DateRange
}

type GetTrips struct {
	DateRange valueobject.DateRange
	City      string
}

type GetDrivers struct {
	DateRange valueobject.DateRange
}

type GetCustomers struct {
	DateRange valueobject.DateRange
}

type GetHeatmap struct {
	City    string
	DateRange valueobject.DateRange
}

type GetForecast struct {
	Metric string
}

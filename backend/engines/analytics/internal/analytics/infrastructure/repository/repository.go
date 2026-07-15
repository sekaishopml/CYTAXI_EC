package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/forecast"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/metrics"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/report"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type ReportRepository interface {
	FindByID(ctx context.Context, id valueobject.ReportID) (*report.AnalyticsReport, error)
	FindByDateRange(ctx context.Context, dr valueobject.DateRange) ([]report.AnalyticsReport, error)
	Save(ctx context.Context, r *report.AnalyticsReport) error
}

type MetricsRepository interface {
	FindByDateRange(ctx context.Context, dr valueobject.DateRange) (*metrics.BusinessMetrics, error)
	Save(ctx context.Context, m *metrics.BusinessMetrics) error
}

type ForecastRepository interface {
	FindByMetric(ctx context.Context, metric string) (*forecast.Forecast, error)
	Save(ctx context.Context, f *forecast.Forecast) error
}

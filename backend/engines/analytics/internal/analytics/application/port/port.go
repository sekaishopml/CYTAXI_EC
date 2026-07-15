package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/business"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/forecast"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/heatmap"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/metrics"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/report"
)

type AnalyticsService interface {
	GenerateReport(ctx context.Context, cmd command.GenerateReport) (*report.AnalyticsReport, error)
	CalculateMetrics(ctx context.Context, cmd command.CalculateMetrics) (*metrics.BusinessMetrics, error)
	GenerateHeatmap(ctx context.Context, cmd command.GenerateHeatmap) (*heatmap.Heatmap, error)
	UpdateDashboard(ctx context.Context, cmd command.UpdateDashboard) (*business.DashboardMetrics, error)
	CreateForecast(ctx context.Context, cmd command.CreateForecast) (*forecast.Forecast, error)
	GetDashboard(ctx context.Context, q query.GetDashboard) (*business.DashboardMetrics, error)
	GetMetrics(ctx context.Context, q query.GetMetrics) (*metrics.BusinessMetrics, error)
	GetRevenue(ctx context.Context, q query.GetRevenue) (*metrics.RevenueMetrics, error)
	GetTrips(ctx context.Context, q query.GetTrips) (*metrics.TripMetrics, error)
	GetDrivers(ctx context.Context, q query.GetDrivers) (*metrics.DriverMetrics, error)
	GetCustomers(ctx context.Context, q query.GetCustomers) (*metrics.CustomerMetrics, error)
	GetHeatmap(ctx context.Context, q query.GetHeatmap) (*heatmap.Heatmap, error)
	GetForecast(ctx context.Context, q query.GetForecast) (*forecast.Forecast, error)
}

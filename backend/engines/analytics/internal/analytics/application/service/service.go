package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/business"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/forecast"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/heatmap"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/metrics"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/report"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/infrastructure/repository"
)

type AnalyticsService struct {
	reportRepo  repository.ReportRepository
	metricsRepo repository.MetricsRepository
	forecastRepo repository.ForecastRepository
	logger      *slog.Logger
}

func NewAnalyticsService(
	reportRepo repository.ReportRepository,
	metricsRepo repository.MetricsRepository,
	forecastRepo repository.ForecastRepository,
	logger *slog.Logger,
) *AnalyticsService {
	return &AnalyticsService{
		reportRepo: reportRepo, metricsRepo: metricsRepo, forecastRepo: forecastRepo, logger: logger,
	}
}

func (s *AnalyticsService) GenerateReport(ctx context.Context, cmd command.GenerateReport) (*report.AnalyticsReport, error) {
	r := report.NewAnalyticsReport(cmd.Title, report.ReportType(cmd.Type), cmd.DateRange)
	r.Granularity = cmd.Granularity
	if err := s.reportRepo.Save(ctx, r); err != nil {
		return nil, fmt.Errorf("save report: %w", err)
	}
	return r, nil
}

func (s *AnalyticsService) CalculateMetrics(ctx context.Context, cmd command.CalculateMetrics) (*metrics.BusinessMetrics, error) {
	m := metrics.NewBusinessMetrics(cmd.DateRange)
	if err := s.metricsRepo.Save(ctx, m); err != nil {
		return nil, fmt.Errorf("save metrics: %w", err)
	}
	return m, nil
}

func (s *AnalyticsService) GenerateHeatmap(ctx context.Context, cmd command.GenerateHeatmap) (*heatmap.Heatmap, error) {
	h := heatmap.NewHeatmap(cmd.City, cmd.DateRange)
	s.logger.Info("heatmap generated", "city", cmd.City)
	return h, nil
}

func (s *AnalyticsService) UpdateDashboard(ctx context.Context, cmd command.UpdateDashboard) (*business.DashboardMetrics, error) {
	d := &business.DashboardMetrics{}
	return d, nil
}

func (s *AnalyticsService) CreateForecast(ctx context.Context, cmd command.CreateForecast) (*forecast.Forecast, error) {
	f := forecast.NewForecast(cmd.Metric, cmd.Method, cmd.Window)
	if err := s.forecastRepo.Save(ctx, f); err != nil {
		return nil, fmt.Errorf("save forecast: %w", err)
	}
	return f, nil
}

func (s *AnalyticsService) GetDashboard(ctx context.Context, q query.GetDashboard) (*business.DashboardMetrics, error) {
	return &business.DashboardMetrics{}, nil
}

func (s *AnalyticsService) GetMetrics(ctx context.Context, q query.GetMetrics) (*metrics.BusinessMetrics, error) {
	dr := valueobject.DateRange(q.DateRange)
	return s.metricsRepo.FindByDateRange(ctx, dr)
}

func (s *AnalyticsService) GetRevenue(ctx context.Context, q query.GetRevenue) (*metrics.RevenueMetrics, error) {
	m, _ := s.metricsRepo.FindByDateRange(ctx, valueobject.DateRange(q.DateRange))
	if m != nil {
		return &m.Revenue, nil
	}
	return &metrics.RevenueMetrics{}, nil
}

func (s *AnalyticsService) GetTrips(ctx context.Context, q query.GetTrips) (*metrics.TripMetrics, error) {
	m, _ := s.metricsRepo.FindByDateRange(ctx, valueobject.DateRange(q.DateRange))
	if m != nil {
		return &m.Trips, nil
	}
	return &metrics.TripMetrics{}, nil
}

func (s *AnalyticsService) GetDrivers(ctx context.Context, q query.GetDrivers) (*metrics.DriverMetrics, error) {
	m, _ := s.metricsRepo.FindByDateRange(ctx, valueobject.DateRange(q.DateRange))
	if m != nil {
		return &m.Drivers, nil
	}
	return &metrics.DriverMetrics{}, nil
}

func (s *AnalyticsService) GetCustomers(ctx context.Context, q query.GetCustomers) (*metrics.CustomerMetrics, error) {
	m, _ := s.metricsRepo.FindByDateRange(ctx, valueobject.DateRange(q.DateRange))
	if m != nil {
		return &m.Customers, nil
	}
	return &metrics.CustomerMetrics{}, nil
}

func (s *AnalyticsService) GetHeatmap(ctx context.Context, q query.GetHeatmap) (*heatmap.Heatmap, error) {
	return heatmap.NewHeatmap(q.City, valueobject.DateRange(q.DateRange)), nil
}

func (s *AnalyticsService) GetForecast(ctx context.Context, q query.GetForecast) (*forecast.Forecast, error) {
	return s.forecastRepo.FindByMetric(ctx, q.Metric)
}

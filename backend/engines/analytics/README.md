# Analytics & Intelligence Engine

Observability, metrics, and business intelligence for the CYTAXI platform.

## Purpose

The Analytics & Intelligence Engine transforms system events into actionable information for analysis, monitoring, and decision-making. It observes, calculates metrics, and generates reports — but never executes business logic.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **AnalyticsReport** — Structured business/operational reports
- **BusinessMetrics** — Trip, revenue, driver, customer, platform metrics
- **OperationalMetrics** — Uptime, response time, error rate, throughput
- **ForecastModel** — Predictive forecasts with confidence intervals

### Entities
- **DashboardMetrics** — Real-time KPIs with charts
- **Heatmap** — Geographic hot zones (cold→peak)
- **Forecast** — Time-series forecasts (moving_avg, exponential, linear)
- **Statistic** — Tagged metric sample
- **BusinessIndicator** — KPI with status (ok/warning/critical)
- **Metric** — Named metric value with trend

### Value Objects
MetricID, ReportID, HeatmapID, ForecastID, DateRange, MetricValue, HeatLevel(4), KPIValue(change/target), TimeGranularity(4)

## Business Metrics (22 predefined)

| Category | Metrics |
|----------|---------|
| **Trip** | total, completed, cancelled, cancel_rate, avg_wait, avg_trip, avg_distance |
| **Revenue** | total_revenue, avg_fare, commission, payouts, platform_fee |
| **Driver** | active, idle, acceptance_rate, avg_rating, top_drivers |
| **Customer** | active_users, new_users, frequent_users, avg_trips_per_user |
| **Platform** | uptime, matching_time, success_rate, promo_coupons_used |

## Future Integrations
Grafana, Prometheus, ClickHouse, BigQuery, Kafka, Flink, Redpanda

## CQRS

**Commands:** GenerateReport, CalculateMetrics, GenerateHeatmap, UpdateDashboard, CreateForecast, RefreshAnalytics

**Queries:** GetDashboard, GetMetrics, GetRevenue, GetTrips, GetDrivers, GetCustomers, GetHeatmap, GetForecast

## Events

| Event | Description |
|-------|-------------|
| `analytics.generated` | Report generated |
| `analytics.metrics_updated` | Metrics calculated |
| `analytics.dashboard_updated` | Dashboard refreshed |
| `analytics.forecast_created` | Forecast created |
| `analytics.heatmap_generated` | Heatmap generated |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /analytics/dashboard | Dashboard |
| GET | /analytics/metrics | Business metrics |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ANALYTICS_PORT` | 8093 | HTTP server port |

## Development

```bash
go run ./cmd/analytics
```

package events

const (
	EventAnalyticsGenerated        = "analytics.generated"
	EventMetricsUpdated           = "analytics.metrics_updated"
	EventDashboardUpdated         = "analytics.dashboard_updated"
	EventForecastCreated          = "analytics.forecast_created"
	EventHeatmapGenerated         = "analytics.heatmap_generated"
	EventBusinessReportGenerated  = "analytics.business_report_generated"
	EventOperationalReportGenerated = "analytics.operational_report_generated"
)

type AnalyticsGeneratedPayload struct {
	ReportID string `json:"report_id"`
	Title    string `json:"title"`
	Type     string `json:"type"`
}

type MetricsUpdatedPayload struct {
	TotalTrips    int     `json:"total_trips"`
	TotalRevenue  float64 `json:"total_revenue"`
	ActiveDrivers int     `json:"active_drivers"`
	ActiveUsers   int     `json:"active_users"`
}

type DashboardUpdatedPayload struct {
	Timestamp string `json:"timestamp"`
	KPICount  int    `json:"kpi_count"`
}

type ForecastCreatedPayload struct {
	Metric     string  `json:"metric"`
	Method     string  `json:"method"`
	Confidence float64 `json:"confidence"`
}

type HeatmapGeneratedPayload struct {
	HeatmapID string `json:"heatmap_id"`
	City      string `json:"city"`
	Zones     int    `json:"zones_count"`
	PeakHours  string `json:"peak_hours,omitempty"`
}

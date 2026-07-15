package report

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type AnalyticsReport struct {
	ID          valueobject.ReportID    `json:"id"`
	Title       string                 `json:"title"`
	Type        ReportType             `json:"type"`
	DateRange   valueobject.DateRange  `json:"date_range"`
	Granularity valueobject.TimeGranularity `json:"granularity"`
	Metrics     []MetricGroup          `json:"metrics"`
	KPIs       []KPIEntry             `json:"kp_is,omitempty"`
	GeneratedAt time.Time              `json:"generated_at"`
}

type ReportType string

const (
	RptBusiness    ReportType = "business"
	RptOperational ReportType = "operational"
	RptCustom      ReportType = "custom"
)

type MetricGroup struct {
	Name    string                   `json:"name"`
	Metrics []valueobject.MetricValue `json:"metrics"`
}

type KPIEntry struct {
	Name  string              `json:"name"`
	Value valueobject.KPIValue `json:"value"`
}

type Report struct {
	ID        valueobject.ReportID    `json:"id"`
	Title     string                 `json:"title"`
	Type      ReportType             `json:"type"`
	Sections  []ReportSection        `json:"sections"`
	CreatedAt time.Time              `json:"created_at"`
}

type ReportSection struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Chart   string `json:"chart_type,omitempty"` // line, bar, pie, heatmap
}

func NewAnalyticsReport(title string, rptType ReportType, dr valueobject.DateRange) *AnalyticsReport {
	return &AnalyticsReport{
		ID:          valueobject.NewReportID(),
		Title:       title,
		Type:        rptType,
		DateRange:   dr,
		Granularity: valueobject.GranDay,
		GeneratedAt: time.Now(),
	}
}

package valueobject

import (
	"fmt"
	"time"
)

type MetricID string
type ReportID string
type HeatmapID string
type ForecastID string

type DateRange struct {
	Start time.Time
	End   time.Time
}

type MetricValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
	Trend string  `json:"trend,omitempty"` // up, down, stable
}

type HeatLevel int

const (
	HeatCold HeatLevel = 0
	HeatWarm HeatLevel = 1
	HeatHot  HeatLevel = 2
	HeatPeak HeatLevel = 3
)

type KPIValue struct {
	Current   float64 `json:"current"`
	Previous  float64 `json:"previous"`
	Change    float64 `json:"change_pct"`
	Target    float64 `json:"target,omitempty"`
}

type TimeGranularity string

const (
	GranHour  TimeGranularity = "hourly"
	GranDay   TimeGranularity = "daily"
	GranWeek  TimeGranularity = "weekly"
	GranMonth TimeGranularity = "monthly"
)

func NewMetricID() MetricID { return MetricID(fmt.Sprintf("met_%d", time.Now().UnixNano())) }
func NewReportID() ReportID { return ReportID(fmt.Sprintf("rpt_%d", time.Now().UnixNano())) }
func NewHeatmapID() HeatmapID { return HeatmapID(fmt.Sprintf("hmp_%d", time.Now().UnixNano())) }
func NewForecastID() ForecastID { return ForecastID(fmt.Sprintf("fcs_%d", time.Now().UnixNano())) }

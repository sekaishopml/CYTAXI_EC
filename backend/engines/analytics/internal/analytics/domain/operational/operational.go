package operational

import "time"

type OperationalMetrics struct {
	ID         string    `json:"id"`
	DateRange  struct {
		Start time.Time
		End   time.Time
	} `json:"date_range"`
	Uptime         float64 `json:"uptime_pct"`
	ResponseTime   float64 `json:"avg_response_ms"`
	ErrorRate      float64 `json:"error_rate_pct"`
	Throughput     float64 `json:"throughput_rps"`
	ActiveServices int     `json:"active_services"`
	Alerts         int     `json:"alerts_count"`
	CalculatedAt   time.Time `json:"calculated_at"`
}

type Statistic struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	Type  string    `json:"type"` // count, avg, sum, distribution
	Value float64   `json:"value"`
	Tags  map[string]string `json:"tags,omitempty"`
	At    time.Time `json:"at"`
}

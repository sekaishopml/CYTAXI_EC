package business

import "time"

type BusinessIndicator struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Current   float64   `json:"current"`
	Previous  float64   `json:"previous"`
	Change    float64   `json:"change_pct"`
	Target    float64   `json:"target"`
	Status    string    `json:"status"` // ok, warning, critical
	UpdatedAt time.Time `json:"updated_at"`
}

type DashboardMetrics struct {
	KPIs     []DashboardKPI `json:"kpis"`
	Charts   []DashboardChart `json:"charts"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type DashboardKPI struct {
	Title string  `json:"title"`
	Value float64 `json:"value"`
	Unit  string  `json:"unit"`
	Icon  string  `json:"icon,omitempty"`
}

type DashboardChart struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Type   string        `json:"type"`
	Labels []string      `json:"labels"`
	Series []ChartSeries `json:"series"`
}

type ChartSeries struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}

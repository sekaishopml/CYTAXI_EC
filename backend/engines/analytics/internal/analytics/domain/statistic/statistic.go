package statistic

import "time"

type Metric struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Value     float64   `json:"value"`
	Unit      string    `json:"unit"`
	Category  string    `json:"category"`
	UpdatedAt time.Time `json:"updated_at"`
}

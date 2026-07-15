package forecast

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/analytics/internal/analytics/domain/valueobject"
)

type ForecastModel struct {
	ID         valueobject.ForecastID  `json:"id"`
	Metric     string                  `json:"metric"`
	Granularity valueobject.TimeGranularity `json:"granularity"`
	DataPoints []ForecastPoint         `json:"data_points"`
	Confidence float64                 `json:"confidence"`
	CreatedAt  time.Time               `json:"created_at"`
}

type ForecastPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	LowBound  float64   `json:"low_bound,omitempty"`
	HighBound float64   `json:"high_bound,omitempty"`
}

type Forecast struct {
	ID        valueobject.ForecastID `json:"id"`
	Metric    string                 `json:"metric"`
	Method    string                 `json:"method"` // moving_avg, exponential, linear
	Window    int                    `json:"window"`
	Result    []ForecastPoint        `json:"result"`
	CreatedAt time.Time              `json:"created_at"`
}

func NewForecast(metric, method string, window int) *Forecast {
	return &Forecast{
		ID:        valueobject.NewForecastID(),
		Metric:    metric,
		Method:    method,
		Window:    window,
		CreatedAt: time.Now(),
	}
}

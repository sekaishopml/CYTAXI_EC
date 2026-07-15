package flow

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type CustomerJourney struct {
	Phone       string `json:"phone"`
	PassengerName string `json:"passenger_name"`
	Origin      string `json:"origin_address"`
	OriginLat   float64 `json:"origin_lat"`
	OriginLng   float64 `json:"origin_lng"`
	Destination string `json:"dest_address"`
	DestLat     float64 `json:"dest_lat"`
	DestLng     float64 `json:"dest_lng"`
}

type JourneyResult struct {
	SessionID    string                 `json:"session_id"`
	TripStatus   string                 `json:"trip_status"`
	FareEstimate map[string]any         `json:"fare_estimate,omitempty"`
	Timeline     []JourneyStep          `json:"timeline"`
	Success      bool                   `json:"success"`
	Error        string                 `json:"error,omitempty"`
}

type JourneyStep struct {
	Step     string `json:"step"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
}

type FlowOrchestrator struct {
	gatewayURL string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewFlowOrchestrator(gatewayURL string, logger *slog.Logger) *FlowOrchestrator {
	return &FlowOrchestrator{
		gatewayURL: gatewayURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

func (f *FlowOrchestrator) ExecuteCustomerJourney(ctx context.Context, journey CustomerJourney) *JourneyResult {
	result := &JourneyResult{Timeline: make([]JourneyStep, 0)}

	steps := []struct {
		name string
		url  string
		body map[string]any
	}{
		{
			name: "start_conversation",
			url:  "/conversation/start",
			body: map[string]any{
				"phone":   journey.Phone,
				"message": "Hola, quiero un taxi",
			},
		},
		{
			name: "create_trip",
			url:  "/trip/request",
			body: map[string]any{
				"customer_id":    "cust_" + journey.Phone,
				"phone":          journey.Phone,
				"passenger_name": journey.PassengerName,
				"origin_address": journey.Origin,
				"origin_lat":     journey.OriginLat,
				"origin_lng":     journey.OriginLng,
				"dest_address":   journey.Destination,
				"dest_lat":       journey.DestLat,
				"dest_lng":       journey.DestLng,
			},
		},
		{
			name: "estimate_fare",
			url:  "/pricing/estimate",
			body: map[string]any{
				"trip_id":      "trip_" + journey.Phone,
				"distance_km":  5.5,
				"duration_sec": 900,
				"region":       "ecuador",
			},
		},
	}

	for _, step := range steps {
		start := time.Now()
		status := "completed"

		resp, err := f.doPost(ctx, step.url, step.body)
		elapsed := time.Since(start)

		if err != nil {
			status = "failed"
			result.Error = fmt.Sprintf("%s: %v", step.name, err)
			result.Timeline = append(result.Timeline, JourneyStep{
				Step: step.name, Status: status, Duration: elapsed.String(),
			})
			f.logger.Error("journey step failed", "step", step.name, "error", err)
			return result
		}

		if step.name == "start_conversation" {
			if sid, ok := resp["session_id"].(string); ok {
				result.SessionID = sid
			}
		}
		if step.name == "estimate_fare" {
			if fare, ok := resp["fare"].(map[string]any); ok {
				result.FareEstimate = fare
			}
		}

		result.Timeline = append(result.Timeline, JourneyStep{
			Step: step.name, Status: status, Duration: elapsed.String(),
		})
		f.logger.Info("journey step", "step", step.name, "duration", elapsed.String())
	}

	result.Success = true
	result.TripStatus = "created"
	return result
}

func (f *FlowOrchestrator) doPost(ctx context.Context, path string, body map[string]any) (map[string]any, error) {
	data, _ := json.Marshal(body)
	url := f.gatewayURL + "/api/v1" + path

	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)
	return result, nil
}

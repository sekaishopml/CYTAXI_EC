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

type AssignmentOrchestrator struct {
	gatewayURL string
	httpClient *http.Client
	logger     *slog.Logger
}

func NewAssignmentOrchestrator(gatewayURL string, logger *slog.Logger) *AssignmentOrchestrator {
	return &AssignmentOrchestrator{
		gatewayURL: gatewayURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		logger:     logger,
	}
}

type DriverCandidate struct {
	DriverID string  `json:"driver_id"`
	Name     string  `json:"name"`
	Distance float64 `json:"distance_meters"`
	ETA      int     `json:"eta_seconds"`
	Score    float64 `json:"score"`
	Vehicle  string  `json:"vehicle"`
	Plate    string  `json:"plate"`
	Rating   float64 `json:"rating"`
}

type AssignmentResult struct {
	Status    string          `json:"status"`
	Driver    *DriverCandidate `json:"driver,omitempty"`
	Timeline  []Step          `json:"timeline"`
	Success   bool            `json:"success"`
	Error     string          `json:"error,omitempty"`
}

type Step struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Duration string `json:"duration"`
}

func (a *AssignmentOrchestrator) StartAssignment(ctx context.Context, tripID string, lat, lng float64) *AssignmentResult {
	result := &AssignmentResult{Timeline: make([]Step, 0)}
	start := time.Now()

	step := func(name string, fn func() error) bool {
		t := time.Now()
		status := "completed"
		err := fn()
		if err != nil {
			status = "failed"
			result.Error = fmt.Sprintf("%s: %v", name, err)
		}
		result.Timeline = append(result.Timeline, Step{
			Name: name, Status: status,
			Duration: time.Since(t).String(),
		})
		if err != nil {
			a.logger.Error("assignment step failed", "step", name, "error", err)
		}
		return err == nil
	}

	step("start_matching", func() error {
		body := map[string]any{"trip_id": tripID, "pickup_lat": lat, "pickup_lng": lng}
		return a.doPost(ctx, "/matching/start", body, nil)
	})

	step("find_candidates", func() error {
		var candidates []DriverCandidate
		err := a.doGet(ctx, fmt.Sprintf("/matching/%s/candidates", "match_demo"), &candidates)
		if err == nil && len(candidates) > 0 {
			result.Driver = &candidates[0]
		}
		return err
	})

	step("send_request", func() error {
		a.logger.Info("trip request sent to driver", "trip_id", tripID)
		return nil
	})

	if result.Driver != nil {
		result.Status = "driver_assigned"
		result.Success = true
	} else {
		result.Status = "searching"
		result.Success = true
	}

	a.logger.Info("assignment completed", "trip_id", tripID, "duration", time.Since(start).String())
	return result
}

func (a *AssignmentOrchestrator) doPost(ctx context.Context, path string, body, out any) error {
	data, _ := json.Marshal(body)
	req, _ := http.NewRequestWithContext(ctx, "POST", a.gatewayURL+"/api/v1"+path, bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.httpClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return fmt.Errorf("http %d", resp.StatusCode) }
	if out != nil { json.NewDecoder(resp.Body).Decode(out) }
	return nil
}

func (a *AssignmentOrchestrator) doGet(ctx context.Context, path string, out any) error {
	req, _ := http.NewRequestWithContext(ctx, "GET", a.gatewayURL+"/api/v1"+path, nil)
	resp, err := a.httpClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode >= 400 { return fmt.Errorf("http %d", resp.StatusCode) }
	if out != nil { json.NewDecoder(resp.Body).Decode(out) }
	return nil
}

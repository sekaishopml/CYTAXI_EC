package testing

import (
	"context"
	"net/http"
	"testing"
	"time"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func TestHealthEndpoint(t *testing.T, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url+"/health", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("health check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("health returned %d", resp.StatusCode)
	}
}

func TestReadyEndpoint(t *testing.T, url string) {
	resp, err := http.Get(url + "/ready")
	if err != nil {
		t.Fatalf("ready check failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("ready returned %d", resp.StatusCode)
	}
}

func TestSmokeEngines(t *testing.T) {
	engines := map[string]string{
		"trip":     "http://localhost:8087",
		"pricing":  "http://localhost:8088",
		"payment":  "http://localhost:8091",
		"customer": "http://localhost:8085",
		"driver":   "http://localhost:8086",
		"matching": "http://localhost:8089",
	}

	for name, url := range engines {
		t.Run(name+"_health", func(t *testing.T) {
			TestHealthEndpoint(t, url)
		})
	}
}

package whatsapp

import "context"

type HealthCheck struct {
	client Client
}

func NewHealthCheck(client Client) *HealthCheck {
	return &HealthCheck{client: client}
}

type HealthResult struct {
	Status   string `json:"status"`
	Session  string `json:"session"`
	Provider string `json:"provider"`
	Error    string `json:"error,omitempty"`
}

func (h *HealthCheck) Check(ctx context.Context) HealthResult {
	status, err := h.client.GetStatus(ctx)
	if err != nil {
		return HealthResult{
			Status: "error",
			Error:  err.Error(),
		}
	}

	return HealthResult{
		Status:  string(status),
		Session: "cytaxi-main",
	}
}

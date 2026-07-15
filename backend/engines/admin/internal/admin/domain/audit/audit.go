package audit

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type AuditLog struct {
	ID         valueobject.AuditID     `json:"id"`
	OperatorID valueobject.OperatorID  `json:"operator_id"`
	Action     valueobject.AuditAction `json:"action"`
	Resource   string                  `json:"resource"`
	ResourceID string                  `json:"resource_id"`
	Changes    map[string]any          `json:"changes,omitempty"`
	IP         string                  `json:"ip,omitempty"`
	CreatedAt  time.Time               `json:"created_at"`
}

func NewAuditLog(operatorID valueobject.OperatorID, action valueobject.AuditAction, resource, resourceID string) *AuditLog {
	return &AuditLog{
		ID:         valueobject.NewAuditID(),
		OperatorID: operatorID,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		CreatedAt:  time.Now(),
	}
}

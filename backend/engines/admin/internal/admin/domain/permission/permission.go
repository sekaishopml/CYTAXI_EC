package permission

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type Permission struct {
	ID          valueobject.PermissionID `json:"id"`
	Name        string                   `json:"name"`
	Resource    string                   `json:"resource"`
	Action      string                   `json:"action"`
	Description string                   `json:"description"`
	CreatedAt   time.Time                `json:"created_at"`
}

type Operator struct {
	ID       valueobject.OperatorID `json:"id"`
	Name     string                 `json:"name"`
	Email    string                 `json:"email"`
	RoleID   valueobject.RoleID    `json:"role_id"`
	Active   bool                  `json:"active"`
	CreatedAt time.Time            `json:"created_at"`
}

func NewPermission(name, resource, action string) *Permission {
	return &Permission{
		ID:        valueobject.NewPermissionID(),
		Name:      name,
		Resource:  resource,
		Action:    action,
		CreatedAt: time.Now(),
	}
}

func NewOperator(name, email string, roleID valueobject.RoleID) *Operator {
	return &Operator{
		ID:        valueobject.NewOperatorID(),
		Name:      name,
		Email:     email,
		RoleID:    roleID,
		Active:    true,
		CreatedAt: time.Now(),
	}
}

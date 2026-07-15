package role

import (
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type Role struct {
	ID          valueobject.RoleID       `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Permissions []valueobject.PermissionID `json:"permissions"`
	System      bool                     `json:"system"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
}

func NewRole(name, description string) *Role {
	now := time.Now()
	return &Role{
		ID:          valueobject.NewRoleID(),
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (r *Role) AddPermission(perm valueobject.PermissionID) {
	r.Permissions = append(r.Permissions, perm)
	r.UpdatedAt = time.Now()
}

func (r *Role) RemovePermission(perm valueobject.PermissionID) {
	for i, p := range r.Permissions {
		if p == perm {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			r.UpdatedAt = time.Now()
			return
		}
	}
}

func (r *Role) HasPermission(perm valueobject.PermissionID) bool {
	for _, p := range r.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

package command

import "github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"

type CreateRole struct {
	Name        string
	Description string
}

type UpdateRole struct {
	RoleID      valueobject.RoleID
	Name        string
	Description string
}

type DeleteRole struct {
	RoleID valueobject.RoleID
}

type AssignPermission struct {
	RoleID       valueobject.RoleID
	PermissionID valueobject.PermissionID
}

type UpdateConfiguration struct {
	Key   valueobject.ConfigurationKey
	Value valueobject.ConfigurationValue
}

type EnableFeature struct {
	FeatureID valueobject.FeatureFlagID
}

type DisableFeature struct {
	FeatureID valueobject.FeatureFlagID
}

type CreateAnnouncement struct {
	Title     string
	Body      string
	Priority  valueobject.Priority
	PublishAt int64
}

type ScheduleMaintenance struct {
	Title       string
	Description string
	StartAt     int64
	EndAt       int64
	Affects     []string
}

type RegisterAudit struct {
	OperatorID valueobject.OperatorID
	Action     valueobject.AuditAction
	Resource   string
	ResourceID string
	Changes    map[string]any
}

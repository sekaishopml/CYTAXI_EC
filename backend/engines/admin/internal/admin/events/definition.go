package events

const (
	EventRoleCreated         = "admin.role_created"
	EventRoleUpdated         = "admin.role_updated"
	EventPermissionAssigned  = "admin.permission_assigned"
	EventConfigurationChanged = "admin.configuration_changed"
	EventFeatureEnabled      = "admin.feature_enabled"
	EventFeatureDisabled     = "admin.feature_disabled"
	EventAnnouncementPublished = "admin.announcement_published"
	EventMaintenanceScheduled = "admin.maintenance_scheduled"
	EventAuditRegistered     = "admin.audit_registered"
)

type RoleCreatedPayload struct {
	RoleID string `json:"role_id"`
	Name   string `json:"name"`
}

type RoleUpdatedPayload struct {
	RoleID  string `json:"role_id"`
	Changes string `json:"changes"`
}

type PermissionAssignedPayload struct {
	RoleID       string `json:"role_id"`
	PermissionID string `json:"permission_id"`
}

type ConfigurationChangedPayload struct {
	Key       string `json:"key"`
	OldValue  string `json:"old_value,omitempty"`
	NewValue  string `json:"new_value"`
}

type FeatureEnabledPayload struct {
	FeatureID string `json:"feature_id"`
	Name      string `json:"name"`
}

type FeatureDisabledPayload struct {
	FeatureID string `json:"feature_id"`
	Name      string `json:"name"`
}

type AnnouncementPublishedPayload struct {
	AnnouncementID string `json:"announcement_id"`
	Title          string `json:"title"`
	Priority       string `json:"priority"`
}

type MaintenanceScheduledPayload struct {
	MaintenanceID string `json:"maintenance_id"`
	Title         string `json:"title"`
	StartAt       string `json:"start_at"`
}

type AuditRegisteredPayload struct {
	AuditID    string `json:"audit_id"`
	Action     string `json:"action"`
	Resource   string `json:"resource"`
}

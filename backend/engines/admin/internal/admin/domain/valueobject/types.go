package valueobject

import (
	"fmt"
	"time"
)

type RoleID string
type PermissionID string
type OperatorID string
type FeatureFlagID string
type AuditID string
type ConfigurationKey string
type ConfigurationValue string
type AnnouncementID string
type MaintenanceID string

type Priority string

const (
	PriLow    Priority = "low"
	PriMedium Priority = "medium"
	PriHigh   Priority = "high"
	PriUrgent Priority = "urgent"
)

type AuditAction string

const (
	AuditCreated AuditAction = "created"
	AuditUpdated AuditAction = "updated"
	AuditDeleted AuditAction = "deleted"
	AuditAccess  AuditAction = "access"
	AuditError   AuditAction = "error"
)

func NewRoleID() RoleID             { return RoleID(fmt.Sprintf("rol_%d", time.Now().UnixNano())) }
func NewPermissionID() PermissionID  { return PermissionID(fmt.Sprintf("prm_%d", time.Now().UnixNano())) }
func NewOperatorID() OperatorID      { return OperatorID(fmt.Sprintf("op_%d", time.Now().UnixNano())) }
func NewFeatureFlagID() FeatureFlagID { return FeatureFlagID(fmt.Sprintf("ff_%d", time.Now().UnixNano())) }
func NewAuditID() AuditID           { return AuditID(fmt.Sprintf("aud_%d", time.Now().UnixNano())) }
func NewAnnouncementID() AnnouncementID { return AnnouncementID(fmt.Sprintf("ann_%d", time.Now().UnixNano())) }
func NewMaintenanceID() MaintenanceID  { return MaintenanceID(fmt.Sprintf("mnt_%d", time.Now().UnixNano())) }

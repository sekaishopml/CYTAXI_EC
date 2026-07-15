# Administration Engine

Platform administration and configuration management for CYTAXI.

## Purpose

The Administration Engine is the administrative brain of CYTAXI. It manages roles, permissions, feature flags, configuration, maintenance windows, audit logs, announcements, and operator management. No other Engine may modify global configuration.

## Architecture

DDD ✓ Clean Architecture ✓ CQRS ✓ Event Driven ✓ Contract First ✓ Zero Trust ✓

## Domain

### Aggregates
- **Role** — Role with assigned permissions (Add/Remove/HasPermission)
- **Configuration** — Global key-value config with scopes (global/city/region)
- **FeatureFlag** — Toggle features with percentage rollout
- **MaintenanceWindow** — Scheduled maintenance affecting engines

### Entities
- **Permission** — Resource + Action permission
- **Operator** — Admin operator with role assignment
- **AuditLog** — Immutable audit trail (created/updated/deleted/access/error)
- **Announcement** — Platform-wide announcements with priority
- **SystemSetting** — System-level parameter

### Value Objects
RoleID, PermissionID, OperatorID, FeatureFlagID, AuditID, ConfigurationKey, ConfigurationValue, AnnouncementID, MaintenanceID, Priority(4), AuditAction(5)

## CQRS

**Commands:** CreateRole, UpdateRole, DeleteRole, AssignPermission, UpdateConfiguration, EnableFeature, DisableFeature, CreateAnnouncement, ScheduleMaintenance, RegisterAudit

**Queries:** GetConfiguration, GetRoles, GetPermissions, GetAuditLogs, GetMaintenance, GetAnnouncements, GetFeatureFlags

## Events

| Event | Description |
|-------|-------------|
| `admin.role_created` | Role created |
| `admin.role_updated` | Role modified |
| `admin.permission_assigned` | Permission assigned to role |
| `admin.configuration_changed` | Configuration value changed |
| `admin.feature_enabled` | Feature toggled on |
| `admin.feature_disabled` | Feature toggled off |
| `admin.announcement_published` | Announcement published |
| `admin.maintenance_scheduled` | Maintenance scheduled |
| `admin.audit_registered` | Audit entry created |

## API

| Method | Path | Description |
|--------|------|-------------|
| GET | /health | Health check |
| GET | /admin/roles | List roles |
| GET | /admin/feature-flags | List feature flags |

## Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `ADMIN_PORT` | 8094 | HTTP server port |

## Development

```bash
go run ./cmd/admin
```

package port

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/command"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/application/query"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/announcement"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/audit"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/configuration"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/featureflag"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/maintenance"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/permission"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/role"
)

type AdminService interface {
	CreateRole(ctx context.Context, cmd command.CreateRole) (*role.Role, error)
	UpdateRole(ctx context.Context, cmd command.UpdateRole) (*role.Role, error)
	DeleteRole(ctx context.Context, cmd command.DeleteRole) error
	AssignPermission(ctx context.Context, cmd command.AssignPermission) (*role.Role, error)
	UpdateConfiguration(ctx context.Context, cmd command.UpdateConfiguration) (*configuration.Configuration, error)
	EnableFeature(ctx context.Context, cmd command.EnableFeature) error
	DisableFeature(ctx context.Context, cmd command.DisableFeature) error
	CreateAnnouncement(ctx context.Context, cmd command.CreateAnnouncement) (*announcement.Announcement, error)
	ScheduleMaintenance(ctx context.Context, cmd command.ScheduleMaintenance) (*maintenance.MaintenanceWindow, error)
	RegisterAudit(ctx context.Context, cmd command.RegisterAudit) (*audit.AuditLog, error)
	GetConfiguration(ctx context.Context, q query.GetConfiguration) (*configuration.Configuration, error)
	GetRoles(ctx context.Context, q query.GetRoles) ([]role.Role, error)
	GetPermissions(ctx context.Context, q query.GetPermissions) ([]permission.Permission, error)
	GetAuditLogs(ctx context.Context, q query.GetAuditLogs) ([]audit.AuditLog, error)
	GetMaintenance(ctx context.Context, q query.GetMaintenance) ([]maintenance.MaintenanceWindow, error)
	GetAnnouncements(ctx context.Context, q query.GetAnnouncements) ([]announcement.Announcement, error)
	GetFeatureFlags(ctx context.Context, q query.GetFeatureFlags) ([]featureflag.FeatureFlag, error)
}

package repository

import (
	"context"

	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/announcement"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/audit"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/configuration"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/featureflag"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/maintenance"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/permission"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/role"
	"github.com/sekaishopml/cytaxi/backend/engines/admin/internal/admin/domain/valueobject"
)

type RoleRepository interface {
	FindByID(ctx context.Context, id valueobject.RoleID) (*role.Role, error)
	FindAll(ctx context.Context) ([]role.Role, error)
	Save(ctx context.Context, r *role.Role) error
	Update(ctx context.Context, r *role.Role) error
	Delete(ctx context.Context, id valueobject.RoleID) error
}

type PermissionRepository interface {
	FindByID(ctx context.Context, id valueobject.PermissionID) (*permission.Permission, error)
	FindAll(ctx context.Context) ([]permission.Permission, error)
	Save(ctx context.Context, p *permission.Permission) error
}

type ConfigurationRepository interface {
	FindByKey(ctx context.Context, key valueobject.ConfigurationKey) (*configuration.Configuration, error)
	FindAll(ctx context.Context) ([]configuration.Configuration, error)
	Save(ctx context.Context, c *configuration.Configuration) error
}

type FeatureFlagRepository interface {
	FindByID(ctx context.Context, id valueobject.FeatureFlagID) (*featureflag.FeatureFlag, error)
	FindAll(ctx context.Context) ([]featureflag.FeatureFlag, error)
	Update(ctx context.Context, f *featureflag.FeatureFlag) error
}

type AuditRepository interface {
	FindAll(ctx context.Context) ([]audit.AuditLog, error)
	Save(ctx context.Context, a *audit.AuditLog) error
}

type MaintenanceRepository interface {
	FindAll(ctx context.Context) ([]maintenance.MaintenanceWindow, error)
	Save(ctx context.Context, m *maintenance.MaintenanceWindow) error
}

type AnnouncementRepository interface {
	FindAll(ctx context.Context) ([]announcement.Announcement, error)
	Save(ctx context.Context, a *announcement.Announcement) error
}

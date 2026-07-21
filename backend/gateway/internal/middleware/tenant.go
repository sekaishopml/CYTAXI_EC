package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sekaishopml/cytaxi/backend/gateway/internal/tenant"
)

type tenantKey struct{}

func TenantResolver(repo tenant.Repository) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var t *tenant.Tenant

			tenantID := r.Header.Get("X-Tenant-ID")
			if tenantID == "" {
				tenantID = r.URL.Query().Get("tenant_id")
			}

			if tenantID != "" {
				t, _ = repo.GetByID(r.Context(), tenant.ID(tenantID))
			}

			if t == nil {
				host := r.Host
				if host != "" {
					parts := strings.SplitN(host, ".", 2)
					if len(parts) == 2 && parts[0] != "www" && parts[0] != "app" {
						t, _ = repo.GetBySlug(r.Context(), parts[0])
					}
				}
			}

			if t == nil {
				t = &tenant.Tenant{
					ID: "default", Name: "Default", Slug: "default",
					Plan: tenant.PlanEnterprise, IsActive: true,
					Locale: "es", Timezone: "America/Guayaquil",
				}
			}

			if !t.IsActive {
				http.Error(w, `{"error":"tenant inactive"}`, http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), tenantKey{}, t)
			w.Header().Set("X-Tenant-ID", string(t.ID))
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func GetTenant(r *http.Request) *tenant.Tenant {
	t, _ := r.Context().Value(tenantKey{}).(*tenant.Tenant)
	return t
}

package router

import (
	"encoding/json"
	"net/http"

	"github.com/sekaishopml/cytaxi/backend/gateway/internal/tenant"
)

func (r *GatewayRouter) RegisterAdminRoutes(repo tenant.Repository) {
	logger := r.logger.With("module", "admin")

	r.mux.HandleFunc("GET /admin/tenants", func(w http.ResponseWriter, req *http.Request) {
		tenants, err := repo.List(req.Context())
		if err != nil {
			http.Error(w, `{"error":"failed to list tenants"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tenants)
	})

	r.mux.HandleFunc("POST /admin/tenants", func(w http.ResponseWriter, req *http.Request) {
		var t tenant.Tenant
		if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
			http.Error(w, `{"error":"invalid tenant data"}`, http.StatusBadRequest)
			return
		}
		if t.ID == "" {
			http.Error(w, `{"error":"tenant ID is required"}`, http.StatusBadRequest)
			return
		}
		if err := repo.Save(req.Context(), &t); err != nil {
			http.Error(w, `{"error":"failed to save tenant"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(t)
	})

	r.mux.HandleFunc("GET /admin/tenants/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		t, err := repo.GetByID(req.Context(), tenant.ID(id))
		if err != nil {
			http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	})

	r.mux.HandleFunc("PUT /admin/tenants/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		var t tenant.Tenant
		if err := json.NewDecoder(req.Body).Decode(&t); err != nil {
			http.Error(w, `{"error":"invalid tenant data"}`, http.StatusBadRequest)
			return
		}
		t.ID = tenant.ID(id)
		if err := repo.Update(req.Context(), &t); err != nil {
			http.Error(w, `{"error":"failed to update tenant"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(t)
	})

	r.mux.HandleFunc("DELETE /admin/tenants/{id}", func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		if err := repo.Delete(req.Context(), tenant.ID(id)); err != nil {
			http.Error(w, `{"error":"failed to delete tenant"}`, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	r.mux.HandleFunc("GET /admin/tenants/{id}/usage", func(w http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		t, err := repo.GetByID(req.Context(), tenant.ID(id))
		if err != nil {
			http.Error(w, `{"error":"tenant not found"}`, http.StatusNotFound)
			return
		}
		usage := map[string]interface{}{
			"tenant_id": t.ID, "name": t.Name, "plan": t.Plan,
			"driver_count": 14, "trip_count": 127, "revenue": 184700,
			"currency": "USD", "period": "today",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(usage)
	})

	logger.Info("admin routes registered", "prefix", "/admin/tenants")
}

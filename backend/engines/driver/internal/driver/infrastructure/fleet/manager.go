package fleet

import (
	"fmt"
	"sync"
	"time"
)

type Fleet struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	OwnerID     string    `json:"owner_id"`
	Description string    `json:"description"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

type FleetVehicle struct {
	ID           string    `json:"id"`
	FleetID      string    `json:"fleet_id"`
	Plate        string    `json:"plate"`
	Brand        string    `json:"brand"`
	Model        string    `json:"model"`
	Year         int       `json:"year"`
	Color        string    `json:"color"`
	Type         string    `json:"type"`
	Status       string    `json:"status"` // active, maintenance, inactive, suspended
	AssignedTo   string    `json:"assigned_to,omitempty"`
	MileageKM    float64   `json:"mileage_km"`
	NextService  *time.Time `json:"next_service_at,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type VehicleAssignment struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	DriverID   string    `json:"driver_id"`
	StartAt    time.Time `json:"start_at"`
	EndAt      *time.Time `json:"end_at,omitempty"`
	Status    string    `json:"status"` // active, ended
	Notes     string    `json:"notes,omitempty"`
}

type MaintenanceRecord struct {
	ID          string    `json:"id"`
	VehicleID   string    `json:"vehicle_id"`
	Type        string    `json:"type"` // oil_change, tire_rotation, inspection, repair, other
	Description string    `json:"description"`
	Status      string    `json:"status"` // scheduled, in_progress, completed
	ScheduledAt time.Time `json:"scheduled_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	Cost        float64   `json:"cost,omitempty"`
	Provider    string    `json:"provider,omitempty"`
}

type Inspection struct {
	ID         string    `json:"id"`
	VehicleID  string    `json:"vehicle_id"`
	Inspector  string    `json:"inspector"`
	Type       string    `json:"type"` // regular, mandatory, incident
	Status     string    `json:"status"` // pending, passed, failed
	Notes      string    `json:"notes,omitempty"`
	Score      float64   `json:"score,omitempty"`
	InspectedAt time.Time `json:"inspected_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

type FleetDashboard struct {
	TotalVehicles    int `json:"total_vehicles"`
	ActiveVehicles   int `json:"active_vehicles"`
	InMaintenance    int `json:"in_maintenance"`
	AssignedDrivers  int `json:"assigned_drivers"`
	PendingInspections int `json:"pending_inspections"`
	FleetUtilization float64 `json:"fleet_utilization_pct"`
}

type Manager struct {
	mu           sync.RWMutex
	fleets       map[string]*Fleet
	vehicles     map[string]*FleetVehicle
	assignments   map[string][]VehicleAssignment
	maintenance  map[string][]MaintenanceRecord
	inspections  map[string][]Inspection
}

func NewManager() *Manager {
	return &Manager{
		fleets:      make(map[string]*Fleet),
		vehicles:    make(map[string]*FleetVehicle),
		assignments: make(map[string][]VehicleAssignment),
		maintenance: make(map[string][]MaintenanceRecord),
		inspections: make(map[string][]Inspection),
	}
}

func (m *Manager) CreateFleet(name, ownerID, desc string) *Fleet {
	f := &Fleet{
		ID: fmt.Sprintf("flt_%d", time.Now().UnixNano()),
		Name: name, OwnerID: ownerID, Description: desc, Active: true, CreatedAt: time.Now(),
	}
	m.mu.Lock()
	m.fleets[f.ID] = f
	m.mu.Unlock()
	return f
}

func (m *Manager) RegisterVehicle(fleetID, plate, brand, model string, year int, vType string) *FleetVehicle {
	v := &FleetVehicle{
		ID: fmt.Sprintf("veh_%d", time.Now().UnixNano()),
		FleetID: fleetID, Plate: plate, Brand: brand, Model: model,
		Year: year, Type: vType, Status: "active", CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	m.mu.Lock()
	m.vehicles[v.ID] = v
	m.mu.Unlock()
	return v
}

func (m *Manager) GetVehicles(fleetID string) []FleetVehicle {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []FleetVehicle
	for _, v := range m.vehicles {
		if fleetID == "" || v.FleetID == fleetID {
			result = append(result, *v)
		}
	}
	return result
}

func (m *Manager) AssignDriver(vehicleID, driverID, notes string) *VehicleAssignment {
	a := &VehicleAssignment{
		ID: fmt.Sprintf("asg_%d", time.Now().UnixNano()),
		VehicleID: vehicleID, DriverID: driverID, StartAt: time.Now(), Status: "active", Notes: notes,
	}
	m.mu.Lock()
	m.assignments[vehicleID] = append(m.assignments[vehicleID], *a)
	if v, ok := m.vehicles[vehicleID]; ok {
		v.AssignedTo = driverID
		v.UpdatedAt = time.Now()
	}
	m.mu.Unlock()
	return a
}

func (m *Manager) ReleaseDriver(vehicleID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := len(m.assignments[vehicleID]) - 1; i >= 0; i-- {
		if m.assignments[vehicleID][i].Status == "active" {
			now := time.Now()
			m.assignments[vehicleID][i].EndAt = &now
			m.assignments[vehicleID][i].Status = "ended"
			if v, ok := m.vehicles[vehicleID]; ok {
				v.AssignedTo = ""
				v.UpdatedAt = now
			}
			return
		}
	}
}

func (m *Manager) ScheduleMaintenance(vehicleID, maintType, desc string) *MaintenanceRecord {
	r := &MaintenanceRecord{
		ID: fmt.Sprintf("mnt_%d", time.Now().UnixNano()),
		VehicleID: vehicleID, Type: maintType, Description: desc,
		Status: "scheduled", ScheduledAt: time.Now(),
	}
	m.mu.Lock()
	m.maintenance[vehicleID] = append(m.maintenance[vehicleID], *r)
	if v, ok := m.vehicles[vehicleID]; ok {
		v.Status = "maintenance"
		v.UpdatedAt = time.Now()
	}
	m.mu.Unlock()
	return r
}

func (m *Manager) CompleteMaintenance(vehicleID, maintenanceID string, cost float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	records := m.maintenance[vehicleID]
	for i := range records {
		if records[i].ID == maintenanceID {
			now := time.Now()
			records[i].Status = "completed"
			records[i].CompletedAt = &now
			records[i].Cost = cost
			m.maintenance[vehicleID] = records
			if v, ok := m.vehicles[vehicleID]; ok {
				v.Status = "active"
				v.UpdatedAt = now
			}
			return nil
		}
	}
	return fmt.Errorf("maintenance record not found")
}

func (m *Manager) GetDashboard() *FleetDashboard {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var active, maintenance, assigned int
	for _, v := range m.vehicles {
		switch v.Status {
		case "active": active++
		case "maintenance": maintenance++
		}
		if v.AssignedTo != "" { assigned++ }
	}

	utilization := 0.0
	if len(m.vehicles) > 0 {
		utilization = float64(active) / float64(len(m.vehicles)) * 100
	}

	return &FleetDashboard{
		TotalVehicles:   len(m.vehicles),
		ActiveVehicles:  active,
		InMaintenance:   maintenance,
		AssignedDrivers: assigned,
		FleetUtilization: utilization,
	}
}

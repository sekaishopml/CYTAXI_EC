package experience

import (
	"fmt"
	"sync"
	"time"
)

type FavoritePlace struct {
	ID         string    `json:"id"`
	CustomerID string    `json:"customer_id"`
	Name       string    `json:"name"`
	Address    string    `json:"address"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Type       string    `json:"type"` // home, work, gym, restaurant, custom
	Icon       string    `json:"icon,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UsedCount  int       `json:"used_count"`
}

type SavedRoute struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Name        string    `json:"name"`
	OriginAddr  string    `json:"origin_address"`
	OriginLat   float64   `json:"origin_lat"`
	OriginLng   float64   `json:"origin_lng"`
	DestAddr    string    `json:"dest_address"`
	DestLat     float64   `json:"dest_lat"`
	DestLng     float64   `json:"dest_lng"`
	Recurring   string    `json:"recurring,omitempty"` // daily, weekly, none
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
}

type CustomerPreference struct {
	CustomerID       string `json:"customer_id"`
	DefaultVehicle   string `json:"default_vehicle"`
	DefaultPayment   string `json:"default_payment"`
	TipPercent       int    `json:"tip_percent"`
	QuietRide        bool   `json:"quiet_ride"`
	TemperaturePref  string `json:"temperature_pref"`
	MusicPref        string `json:"music_pref,omitempty"`
	Language         string `json:"language"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type LoyaltyAccount struct {
	CustomerID  string    `json:"customer_id"`
	Points      int       `json:"points"`
	Level       string    `json:"level"` // bronze, silver, gold, platinum
	TripsTotal  int       `json:"trips_total"`
	SpentTotal  float64   `json:"spent_total"`
	JoinedAt    time.Time `json:"joined_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LoyaltyReward struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PointsCost  int    `json:"points_cost"`
	Type        string `json:"type"` // discount, free_ride, upgrade
	Value       float64 `json:"value"`
}

type SupportTicket struct {
	ID          string    `json:"id"`
	CustomerID  string    `json:"customer_id"`
	Subject     string    `json:"subject"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // trip, payment, driver, app, other
	Priority    string    `json:"priority"` // low, medium, high, urgent
	Status      string    `json:"status"`   // open, in_progress, resolved, closed
	Resolution  string    `json:"resolution,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type NotificationPref struct {
	CustomerID  string `json:"customer_id"`
	TripUpdates bool   `json:"trip_updates"`
	PromoOffers bool   `json:"promo_offers"`
	Receipts    bool   `json:"receipts"`
	DriverInfo  bool   `json:"driver_info"`
	Channel     string `json:"channel"` // push, email, sms, whatsapp
}

type Manager struct {
	mu           sync.RWMutex
	favorites    map[string][]FavoritePlace
	savedRoutes  map[string][]SavedRoute
	preferences  map[string]*CustomerPreference
	loyalty      map[string]*LoyaltyAccount
	tickets      map[string]*SupportTicket
	notifPrefs   map[string]*NotificationPref
}

func NewManager() *Manager {
	return &Manager{
		favorites:   make(map[string][]FavoritePlace),
		savedRoutes: make(map[string][]SavedRoute),
		preferences: make(map[string]*CustomerPreference),
		loyalty:     make(map[string]*LoyaltyAccount),
		tickets:     make(map[string]*SupportTicket),
		notifPrefs:  make(map[string]*NotificationPref),
	}
}

func (m *Manager) AddFavorite(customerID string, f FavoritePlace) *FavoritePlace {
	f.ID = fmt.Sprintf("fav_%d", time.Now().UnixNano())
	f.CustomerID = customerID
	f.CreatedAt = time.Now()
	m.mu.Lock()
	m.favorites[customerID] = append(m.favorites[customerID], f)
	m.mu.Unlock()
	return &f
}

func (m *Manager) GetFavorites(customerID string) []FavoritePlace {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.favorites[customerID]
}

func (m *Manager) SaveRoute(customerID string, r SavedRoute) *SavedRoute {
	r.ID = fmt.Sprintf("rte_%d", time.Now().UnixNano())
	r.CustomerID = customerID
	r.CreatedAt = time.Now()
	m.mu.Lock()
	m.savedRoutes[customerID] = append(m.savedRoutes[customerID], r)
	m.mu.Unlock()
	return &r
}

func (m *Manager) GetSavedRoutes(customerID string) []SavedRoute {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.savedRoutes[customerID]
}

func (m *Manager) UpdatePreferences(customerID string, prefs CustomerPreference) *CustomerPreference {
	prefs.CustomerID = customerID
	prefs.UpdatedAt = time.Now()
	m.mu.Lock()
	m.preferences[customerID] = &prefs
	m.mu.Unlock()
	return &prefs
}

func (m *Manager) GetPreferences(customerID string) *CustomerPreference {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.preferences[customerID]
	if !ok {
		return &CustomerPreference{
			CustomerID: customerID, DefaultVehicle: "standard",
			DefaultPayment: "card", Language: "es", TipPercent: 10,
		}
	}
	return p
}

func (m *Manager) GetLoyalty(customerID string) *LoyaltyAccount {
	m.mu.RLock()
	defer m.mu.RUnlock()
	acc, ok := m.loyalty[customerID]
	if !ok {
		return &LoyaltyAccount{
			CustomerID: customerID, Points: 0, Level: "bronze",
			JoinedAt: time.Now(), UpdatedAt: time.Now(),
		}
	}
	return acc
}

func (m *Manager) EarnPoints(customerID string, tripAmount float64) *LoyaltyAccount {
	m.mu.Lock()
	defer m.mu.Unlock()
	acc, ok := m.loyalty[customerID]
	if !ok {
		acc = &LoyaltyAccount{
			CustomerID: customerID, Points: 0, Level: "bronze",
			JoinedAt: time.Now(),
		}
		m.loyalty[customerID] = acc
	}
	points := int(tripAmount * 2) // $1 = 2 points
	acc.Points += points
	acc.TripsTotal++
	acc.SpentTotal += tripAmount
	acc.Level = calculateLevel(acc.Points)
	acc.UpdatedAt = time.Now()
	return acc
}

func calculateLevel(points int) string {
	switch {
	case points >= 10000: return "platinum"
	case points >= 5000: return "gold"
	case points >= 2000: return "silver"
	default: return "bronze"
	}
}

func (m *Manager) CreateTicket(customerID, subject, description, category, priority string) *SupportTicket {
	t := &SupportTicket{
		ID: fmt.Sprintf("tkt_%d", time.Now().UnixNano()),
		CustomerID: customerID, Subject: subject, Description: description,
		Category: category, Priority: priority, Status: "open",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	m.mu.Lock()
	m.tickets[t.ID] = t
	m.mu.Unlock()
	return t
}

func (m *Manager) GetTickets(customerID string) []*SupportTicket {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*SupportTicket
	for _, t := range m.tickets {
		if t.CustomerID == customerID {
			result = append(result, t)
		}
	}
	return result
}

func (m *Manager) ResolveTicket(ticketID, resolution string) (*SupportTicket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tickets[ticketID]
	if !ok { return nil, fmt.Errorf("ticket not found") }
	now := time.Now()
	t.Status = "resolved"
	t.Resolution = resolution
	t.ResolvedAt = &now
	t.UpdatedAt = now
	return t, nil
}

func (m *Manager) UpdateNotificationPrefs(customerID string, prefs NotificationPref) *NotificationPref {
	prefs.CustomerID = customerID
	m.mu.Lock()
	m.notifPrefs[customerID] = &prefs
	m.mu.Unlock()
	return &prefs
}

func (m *Manager) GetNotificationPrefs(customerID string) *NotificationPref {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.notifPrefs[customerID]
	if !ok {
		return &NotificationPref{
			CustomerID: customerID, TripUpdates: true, Receipts: true,
			Channel: "push",
		}
	}
	return p
}

func (m *Manager) RedeemPoints(customerID string, points int) (*LoyaltyAccount, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	acc, ok := m.loyalty[customerID]
	if !ok { return nil, fmt.Errorf("loyalty account not found") }
	if acc.Points < points { return nil, fmt.Errorf("insufficient points") }
	acc.Points -= points
	acc.UpdatedAt = time.Now()
	return acc, nil
}

var rewards = []LoyaltyReward{
	{ID: "rwd_1", Name: "Free Ride (up to $10)", PointsCost: 500, Type: "free_ride", Value: 10},
	{ID: "rwd_2", Name: "50% Discount", PointsCost: 300, Type: "discount", Value: 50},
	{ID: "rwd_3", Name: "Vehicle Upgrade", PointsCost: 1000, Type: "upgrade", Value: 0},
}

func GetRewards() []LoyaltyReward { return rewards }

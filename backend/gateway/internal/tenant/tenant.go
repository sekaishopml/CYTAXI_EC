package tenant

type ID string

type Plan string

const (
	PlanFree       Plan = "free"
	PlanStarter    Plan = "starter"
	PlanBusiness   Plan = "business"
	PlanEnterprise Plan = "enterprise"
)

type Tenant struct {
	ID          ID       `json:"id"`
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Plan        Plan     `json:"plan"`
	IsActive    bool     `json:"is_active"`
	MaxDrivers  int      `json:"max_drivers"`
	MaxVehicles int      `json:"max_vehicles"`
	Locale      string   `json:"locale"`
	Timezone    string   `json:"timezone"`
	Domain      string   `json:"domain"`
	Branding    Branding `json:"branding"`
	Features    []string `json:"features"`
	CreatedAt   int64    `json:"created_at"`
	UpdatedAt   int64    `json:"updated_at"`
}

type Branding struct {
	PrimaryColor   string `json:"primary_color"`
	SecondaryColor string `json:"secondary_color"`
	LogoURL        string `json:"logo_url"`
	FaviconURL     string `json:"favicon_url"`
	AppName        string `json:"app_name"`
	TermsURL       string `json:"terms_url"`
	PrivacyURL     string `json:"privacy_url"`
}

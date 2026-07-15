package profile

import "time"

type Profile struct {
	CustomerID  string    `json:"customer_id"`
	FullName    string    `json:"full_name"`
	Phone       string    `json:"phone"`
	Email       string    `json:"email,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Language    string    `json:"language"`
	Timezone    string    `json:"timezone"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewProfile(customerID, phone string) *Profile {
	return &Profile{
		CustomerID: customerID,
		Phone:      phone,
		Language:   "es",
		Timezone:   "America/Guayaquil",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (p *Profile) UpdateName(name string) {
	p.FullName = name
	p.UpdatedAt = time.Now()
}

func (p *Profile) UpdateEmail(email string) {
	p.Email = email
	p.UpdatedAt = time.Now()
}

func (p *Profile) UpdateLanguage(lang string) {
	p.Language = lang
	p.UpdatedAt = time.Now()
}

func (p *Profile) UpdateTimezone(tz string) {
	p.Timezone = tz
	p.UpdatedAt = time.Now()
}

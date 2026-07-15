package recipient

import "github.com/sekaishopml/cytaxi/backend/engines/notification/internal/notification/domain/valueobject"

type Recipient struct {
	ID    valueobject.RecipientID    `json:"id"`
	Phone string                     `json:"phone,omitempty"`
	Email string                     `json:"email,omitempty"`
	Name  string                     `json:"name"`
	Locale valueobject.Locale        `json:"locale"`
	Devices []Device                 `json:"devices,omitempty"`
}

type Device struct {
	ID     string `json:"id"`
	Type   string `json:"type"` // ios, android, web
	Token  string `json:"token"`
	Active bool   `json:"active"`
}

func (r *Recipient) ActiveDevices() []Device {
	var active []Device
	for _, d := range r.Devices {
		if d.Active {
			active = append(active, d)
		}
	}
	return active
}

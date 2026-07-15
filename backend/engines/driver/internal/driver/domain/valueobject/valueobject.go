package valueobject

import (
	"fmt"
	"regexp"
	"strings"
)

type DriverID string
type VehicleID string
type LicenseNumber string
type PlateNumber string
type Phone string
type Email string
type Coordinates struct {
	Lat float64
	Lng float64
}

func NewPhone(p string) (Phone, error) {
	p = strings.TrimSpace(p)
	if len(p) < 10 {
		return "", fmt.Errorf("invalid phone: %s", p)
	}
	return Phone(p), nil
}

func NewEmail(e string) (Email, error) {
	e = strings.TrimSpace(strings.ToLower(e))
	if !strings.Contains(e, "@") || !strings.Contains(e, ".") {
		return "", fmt.Errorf("invalid email: %s", e)
	}
	return Email(e), nil
}

var plateRegex = regexp.MustCompile(`^[A-Z]{3}-?\d{3,4}$`)

func NewPlateNumber(p string) (PlateNumber, error) {
	p = strings.TrimSpace(strings.ToUpper(p))
	if !plateRegex.MatchString(p) {
		return "", fmt.Errorf("invalid plate: %s", p)
	}
	return PlateNumber(p), nil
}

func (p Phone) String() string    { return string(p) }
func (e Email) String() string    { return string(e) }
func (l LicenseNumber) String() string { return string(l) }
func (p PlateNumber) String() string { return string(p) }
func (d DriverID) String() string { return string(d) }
func (v VehicleID) String() string { return string(v) }
func (c Coordinates) IsZero() bool { return c.Lat == 0 && c.Lng == 0 }

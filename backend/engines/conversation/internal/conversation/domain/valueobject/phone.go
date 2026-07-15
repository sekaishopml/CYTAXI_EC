package valueobject

import "fmt"

type PhoneNumber string

func NewPhoneNumber(number string) (PhoneNumber, error) {
	if len(number) < 10 {
		return "", fmt.Errorf("invalid phone number: %s", number)
	}
	return PhoneNumber(number), nil
}

func (p PhoneNumber) String() string {
	return string(p)
}

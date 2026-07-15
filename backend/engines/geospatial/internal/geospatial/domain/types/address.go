package types

type Address struct {
	FormattedAddress string      `json:"formatted_address"`
	Street           string      `json:"street"`
	City             string      `json:"city"`
	State            string      `json:"state"`
	Country          string      `json:"country"`
	PostalCode       string      `json:"postal_code"`
	Coordinates      Coordinates `json:"coordinates"`
}

type AddressComponent struct {
	LongName  string   `json:"long_name"`
	ShortName string   `json:"short_name"`
	Types     []string `json:"types"`
}

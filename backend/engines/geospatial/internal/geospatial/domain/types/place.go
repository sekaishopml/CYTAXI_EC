package types

type Place struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Coordinates Coordinates `json:"coordinates"`
	Types       []string    `json:"types"`
	Rating      float64     `json:"rating,omitempty"`
	PhoneNumber string      `json:"phone_number,omitempty"`
}

type PlaceSearchRequest struct {
	Query     string      `json:"query"`
	Location  Coordinates `json:"location,omitempty"`
	Radius    int         `json:"radius,omitempty"` // meters
	Types     []string    `json:"types,omitempty"`
	MaxResult int         `json:"max_results,omitempty"`
}

type PlaceSearchResult struct {
	Places     []Place `json:"places"`
	NextPage   string  `json:"next_page,omitempty"`
	TotalCount int     `json:"total_count"`
}

type AutocompleteRequest struct {
	Input    string      `json:"input"`
	Location Coordinates `json:"location,omitempty"`
	Radius   int         `json:"radius,omitempty"`
	Types    []string    `json:"types,omitempty"`
}

type AutocompletePrediction struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Types       []string `json:"types"`
}

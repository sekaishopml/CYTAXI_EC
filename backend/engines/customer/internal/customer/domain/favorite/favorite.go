package favorite

import (
	"fmt"
	"time"
)

type FavoritePlaceID string

type FavoritePlace struct {
	ID          FavoritePlaceID `json:"id"`
	CustomerID  string          `json:"customer_id"`
	Name        string          `json:"name"`
	Address     string          `json:"address"`
	Lat         float64         `json:"lat"`
	Lng         float64         `json:"lng"`
	Category    string          `json:"category,omitempty"` // home, work, gym, etc.
	Frequencies int             `json:"frequencies"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func NewFavoritePlace(customerID, name, address string, lat, lng float64) *FavoritePlace {
	now := time.Now()
	return &FavoritePlace{
		ID:          FavoritePlaceID(fmt.Sprintf("fav_%d", now.UnixNano())),
		CustomerID:  customerID,
		Name:        name,
		Address:     address,
		Lat:         lat,
		Lng:         lng,
		Frequencies: 1,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type FavoritePlaceCategory string

const (
	CategoryHome  FavoritePlaceCategory = "home"
	CategoryWork  FavoritePlaceCategory = "work"
	CategoryGym   FavoritePlaceCategory = "gym"
	CategoryHotel FavoritePlaceCategory = "hotel"
	CategoryOther FavoritePlaceCategory = "other"
)

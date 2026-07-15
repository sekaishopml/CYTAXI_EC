package valueobject

import (
	"fmt"
	"time"
)

type MatchingID string
type DriverID string
type TripID string
type SessionID string

type Distance struct {
	Meters float64
	KM     float64
}

type ETA struct {
	Seconds int
	Minutes float64
}

type MatchingScore float64

type Priority int

type AvailabilityStatus string

const (
	AvailAvailable AvailabilityStatus = "available"
	AvailBusy      AvailabilityStatus = "busy"
	AvailOnTrip    AvailabilityStatus = "on_trip"
	AvailOffline   AvailabilityStatus = "offline"
)

type CandidateRank int

func NewMatchingID() MatchingID { return MatchingID(fmt.Sprintf("match_%d", time.Now().UnixNano())) }

func NewDistance(meters float64) Distance {
	return Distance{Meters: meters, KM: meters / 1000}
}

func NewETA(seconds int) ETA {
	return ETA{Seconds: seconds, Minutes: float64(seconds) / 60}
}

func (s MatchingScore) Float() float64 { return float64(s) }

type DriverSnapshot struct {
	DriverID   DriverID
	Location   struct{ Lat, Lng float64 }
	Distance   Distance
	ETA        ETA
	Score      MatchingScore
	Rating     float64
	TripsCount int
}

type Coordinates struct {
	Lat float64
	Lng float64
}

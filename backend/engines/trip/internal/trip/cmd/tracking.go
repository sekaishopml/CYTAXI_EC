package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/sekaishopml/cytaxi/backend/engines/trip/internal/trip/application/port"
)

type DriverLocation struct {
	DriverID string  `json:"driver_id"`
	TripID   string  `json:"trip_id"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Speed    float64 `json:"speed"`
	Heading  float64 `json:"heading"`
	ETA      int     `json:"eta_seconds"`
	Distance float64 `json:"distance_km"`
	UpdatedAt string `json:"updated_at"`
}

type TrackingServer struct {
	service  port.TripService
	logger   *slog.Logger
	locations sync.Map // tripID → current driver location
	clients   sync.Map // tripID → []chan TrackingUpdate
}

type TrackingUpdate struct {
	Type      string  `json:"type"`
	TripID    string  `json:"trip_id"`
	Status    string  `json:"status"`
	Driver    *TrackingDriver `json:"driver,omitempty"`
	ETA       int     `json:"eta_seconds,omitempty"`
	Timestamp string  `json:"timestamp"`
}

type TrackingDriver struct {
	ID      string  `json:"id"`
	Name    string  `json:"name"`
	Vehicle string  `json:"vehicle"`
	Plate   string  `json:"plate"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
	Rating  float64 `json:"rating"`
}

var simLocation = struct {
	lat, lng float64
	step     int
}{-0.180653, -78.467838, 0}

func NewTrackingServer(service port.TripService, logger *slog.Logger) *TrackingServer {
	return &TrackingServer{service: service, logger: logger}
}

func (ts *TrackingServer) Subscribe(tripID string) <-chan TrackingUpdate {
	ch := make(chan TrackingUpdate, 10)
	val, _ := ts.clients.LoadOrStore(tripID, []chan TrackingUpdate{})
	chans := val.([]chan TrackingUpdate)
	ts.clients.Store(tripID, append(chans, ch))
	return ch
}

func (ts *TrackingServer) broadcast(tripID string, update TrackingUpdate) {
	val, ok := ts.clients.Load(tripID)
	if !ok { return }
	for _, ch := range val.([]chan TrackingUpdate) {
		select {
		case ch <- update:
		default:
		}
	}
}

func (ts *TrackingServer) HandleStartTrip(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID   string `json:"trip_id"`
		DriverID string `json:"driver_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	ts.logger.Info("trip started", "trip_id", req.TripID, "driver", req.DriverID)

	loc := DriverLocation{
		DriverID: req.DriverID, TripID: req.TripID,
		Lat: simLocation.lat, Lng: simLocation.lng,
		ETA: 600, Distance: 5.5,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	ts.locations.Store(req.TripID, loc)

	ts.broadcast(req.TripID, TrackingUpdate{
		Type: "trip_started", TripID: req.TripID, Status: "in_progress",
		Driver: &TrackingDriver{
			ID: req.DriverID, Name: "Carlos M.", Vehicle: "Toyota Corolla",
			Plate: "ABC-1234", Lat: loc.Lat, Lng: loc.Lng, Rating: 4.8,
		},
		ETA: loc.ETA, Timestamp: loc.UpdatedAt,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"status": "started", "trip_id": req.TripID,
		"driver": TrackingDriver{ID: req.DriverID, Name: "Carlos M.", Vehicle: "Toyota Corolla", Plate: "ABC-1234", Rating: 4.8},
	})
}

func (ts *TrackingServer) HandleUpdateLocation(w http.ResponseWriter, r *http.Request) {
	var req DriverLocation
	json.NewDecoder(r.Body).Decode(&req)

	// Simulate movement toward destination
	simLocation.step++
	simLocation.lat += 0.0001
	simLocation.lng += 0.00015

	loc := DriverLocation{
		DriverID: req.DriverID, TripID: req.TripID,
		Lat: simLocation.lat, Lng: simLocation.lng,
		Speed: 40, Heading: 225,
		ETA: 600 - simLocation.step*2,
		Distance: 5.5 - float64(simLocation.step)*0.1,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	ts.locations.Store(req.TripID, loc)

	ts.broadcast(req.TripID, TrackingUpdate{
		Type: "location_update", TripID: req.TripID, Status: "in_progress",
		Driver: &TrackingDriver{
			ID: req.DriverID, Name: "Carlos M.", Vehicle: "Toyota Corolla",
			Plate: "ABC-1234", Lat: loc.Lat, Lng: loc.Lng, Rating: 4.8,
		},
		ETA: loc.ETA, Timestamp: loc.UpdatedAt,
	})

	writeJSON(w, http.StatusOK, map[string]any{"status": "updated", "eta": loc.ETA, "distance_km": loc.Distance})
}

func (ts *TrackingServer) HandleGetLocation(w http.ResponseWriter, r *http.Request) {
	tripID := r.PathValue("trip_id")
	loc, _ := ts.locations.Load(tripID)
	writeJSON(w, http.StatusOK, map[string]any{"location": loc})
}

func (ts *TrackingServer) HandleFinishTrip(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TripID string `json:"trip_id"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	ts.logger.Info("trip finished", "trip_id", req.TripID)

	ts.broadcast(req.TripID, TrackingUpdate{
		Type: "trip_completed", TripID: req.TripID, Status: "completed",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	ts.locations.Delete(req.TripID)
	ts.clients.Delete(req.TripID)

	writeJSON(w, http.StatusOK, map[string]any{"status": "completed", "trip_id": req.TripID})
}

func (ts *TrackingServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	tripID := r.URL.Query().Get("trip_id")
	if tripID == "" {
		http.Error(w, "trip_id required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok { return }

	ch := ts.Subscribe(tripID)
	defer func() {
		val, _ := ts.clients.Load(tripID)
		if val != nil {
			chans := val.([]chan TrackingUpdate)
			for i, c := range chans {
				if c == ch {
					ts.clients.Store(tripID, append(chans[:i], chans[i+1:]...))
					break
				}
			}
		}
	}()

	// Send initial state
	if loc, ok := ts.locations.Load(tripID); ok {
		l := loc.(DriverLocation)
		data, _ := json.Marshal(TrackingUpdate{
			Type: "initial", TripID: tripID, Status: "in_progress",
			Driver: &TrackingDriver{ID: l.DriverID, Name: "Carlos M.", Vehicle: "Toyota Corolla", Plate: "ABC-1234", Lat: l.Lat, Lng: l.Lng, Rating: 4.8},
			ETA: l.ETA, Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		w.Write([]byte("data: " + string(data) + "\n\n"))
		flusher.Flush()
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case update := <-ch:
			data, _ := json.Marshal(update)
			w.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()
		case <-ticker.C:
			// Simulate movement
			simLocation.step++
			simLocation.lat += 0.0001
			simLocation.lng += 0.00015
			eta := max(0, 600-simLocation.step*2)
			distance := max(0, 5.5-float64(simLocation.step)*0.1)

			update := TrackingUpdate{
				Type: "location_update", TripID: tripID, Status: "in_progress",
				Driver: &TrackingDriver{ID: "drv_1000", Name: "Carlos M.", Vehicle: "Toyota Corolla", Plate: "ABC-1234", Lat: simLocation.lat, Lng: simLocation.lng, Rating: 4.8},
				ETA: eta, Timestamp: time.Now().UTC().Format(time.RFC3339),
			}
			loc := DriverLocation{DriverID: "drv_1000", TripID: tripID, Lat: simLocation.lat, Lng: simLocation.lng, ETA: eta, Distance: distance}
			ts.locations.Store(tripID, loc)

			data, _ := json.Marshal(update)
			w.Write([]byte("data: " + string(data) + "\n\n"))
			flusher.Flush()

			if distance <= 0.1 {
				finish := TrackingUpdate{Type: "trip_completed", TripID: tripID, Status: "completed", Timestamp: time.Now().UTC().Format(time.RFC3339)}
				d, _ := json.Marshal(finish)
				w.Write([]byte("data: " + string(d) + "\n\n"))
				flusher.Flush()
				return
			}
		case <-r.Context().Done():
			return
		}
	}
}

func (ts *TrackingServer) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "service": "trip-tracking"})
}

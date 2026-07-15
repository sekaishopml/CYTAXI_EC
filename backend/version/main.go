package main

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
	"time"
)

var (
	Version   = "1.0.0-rc2"
	BuildDate = "2026-07-15"
	CommitSHA = "unknown"
)

type BuildInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"build_date"`
	Commit    string `json:"commit"`
	GoVersion string `json:"go_version"`
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Uptime    string `json:"uptime"`
}

var startTime = time.Now()

func HandleVersion(w http.ResponseWriter, r *http.Request) {
	info := BuildInfo{
		Version:   Version,
		BuildDate: BuildDate,
		Commit:    CommitSHA,
		GoVersion: runtime.Version(),
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Uptime:    time.Since(startTime).String(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func HandleBuild(w http.ResponseWriter, r *http.Request) {
	info := map[string]string{
		"version":    Version,
		"build_date": BuildDate,
		"commit":     CommitSHA,
		"go_version": runtime.Version(),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /version", HandleVersion)
	mux.HandleFunc("GET /build", HandleBuild)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok","version":"` + Version + `"}`))
	})

	port := os.Getenv("PORT")
	if port == "" { port = "8000" }

	http.ListenAndServe(":"+port, mux)
}

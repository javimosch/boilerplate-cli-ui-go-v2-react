package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

//go:embed ui/*
var uiFiles embed.FS

var (
	server *http.Server
	mu     sync.Mutex
)

type Status struct {
	Status    string `json:"status"`
	Port      int    `json:"port"`
	Uptime    string `json:"uptime"`
	Version   string `json:"version"`
	StartTime time.Time `json:"start_time"`
}

var serverStatus Status

func startServer(port int) {
	serverStatus = Status{
		Status:    "running",
		Port:      port,
		Version:   Version,
		StartTime: time.Now(),
	}

	mux := http.NewServeMux()

	// API endpoints (agent-first: JSON responses)
	mux.HandleFunc("/api/status", handleStatusAPI)
	mux.HandleFunc("/api/health", handleHealthAPI)

	// Static UI files from embedded filesystem
	uiSub, _ := fs.Sub(uiFiles, "ui")
	fileServer := http.FileServer(http.FS(uiSub))
	mux.Handle("/", fileServer)

	server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	log.Printf("Server starting on http://localhost:%d", port)
	log.Printf("UI available at http://localhost:%d/", port)
	log.Printf("API available at http://localhost:%d/api/status", port)
	log.Printf("Press Ctrl+C to stop")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

func handleStatusAPI(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	serverStatus.Uptime = time.Since(serverStatus.StartTime).String()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serverStatus)
}

func handleHealthAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"version": Version,
	})
}

// contentType returns the MIME type for a file path
func contentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".json":
		return "application/json"
	case ".svg":
		return "image/svg+xml"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".ico":
		return "image/x-icon"
	default:
		return "application/octet-stream"
	}
}

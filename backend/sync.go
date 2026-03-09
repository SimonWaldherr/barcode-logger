// Barcode-Logger – Minimalistic Sync Backend (Go)
//
// Receives scan entries as HTTP POST (application/json) from the Barcode-Logger
// PWA webhook and appends them to a newline-delimited JSON log file.
//
// Build
// -----
//
//	go build -o sync-backend sync.go
//
// Run
// ---
//
//	./sync-backend -addr :8080 -log scans.ndjson
//
// Then point the Barcode-Logger Webhook URL to:
//
//	http://<your-host>:8080/scan
//
// Security note
// -------------
// This is a minimal example. For production use, run behind a reverse proxy
// (nginx / caddy) that handles TLS, and add authentication (e.g. a shared
// secret checked in the Authorization header).
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// ScanEntry mirrors the JSON payload sent by the Barcode-Logger PWA.
// Unknown fields are preserved via the Extra map.
type ScanEntry struct {
	ID         any     `json:"id,omitempty"`
	Code       string  `json:"code"`
	Format     string  `json:"format,omitempty"`
	Mode       string  `json:"mode,omitempty"`
	Qty        any     `json:"qty,omitempty"`
	Shelf      string  `json:"shelf,omitempty"`
	Recipient  string  `json:"recipient,omitempty"`
	Reference  string  `json:"reference,omitempty"`
	WorkType   string  `json:"workType,omitempty"`
	Technician string  `json:"technician,omitempty"`
	Status     string  `json:"status,omitempty"`
	Notes      string  `json:"notes,omitempty"`
	Ts         int64   `json:"ts,omitempty"`
	Lat        float64 `json:"lat,omitempty"`
	Lon        float64 `json:"lon,omitempty"`
	ReceivedAt string  `json:"received_at,omitempty"`
}

var (
	logPath string
	mu      sync.Mutex
)

func corsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	corsHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	fmt.Fprintf(w, `{"error":%q}`, msg)
}

func handleScan(w http.ResponseWriter, r *http.Request) {
	corsHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MiB limit
	if err != nil || len(body) == 0 {
		jsonError(w, http.StatusBadRequest, "empty or unreadable body")
		return
	}

	var entry ScanEntry
	if err := json.Unmarshal(body, &entry); err != nil {
		jsonError(w, http.StatusUnprocessableEntity, "invalid JSON: "+err.Error())
		return
	}
	if entry.Code == "" {
		jsonError(w, http.StatusUnprocessableEntity, `missing required field "code"`)
		return
	}

	entry.ReceivedAt = time.Now().UTC().Format(time.RFC3339)

	line, err := json.Marshal(entry)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "marshal error")
		return
	}

	mu.Lock()
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err == nil {
		_, err = fmt.Fprintf(f, "%s\n", line)
		f.Close()
	}
	mu.Unlock()

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "log write error: "+err.Error())
		return
	}

	log.Printf("scan received: code=%q mode=%q", entry.Code, entry.Mode)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"ok":true,"code":%q,"received_at":%q}`, entry.Code, entry.ReceivedAt)
}

func main() {
	addr    := flag.String("addr", ":8080", "listen address (e.g. :8080 or 0.0.0.0:9000)")
	logFile := flag.String("log", "scans.ndjson", "NDJSON log file path")
	flag.Parse()

	logPath = *logFile

	http.HandleFunc("/scan", handleScan)
	http.HandleFunc("/", handleScan) // also accept root path

	log.Printf("Barcode-Logger sync backend listening on %s", *addr)
	log.Printf("Logging to: %s", logPath)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal(err)
	}
}

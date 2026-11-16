package api

import (
	"encoding/json"
	"net/http"

	"github.com/asim/reminder/model"
)

// CapabilitiesHandler exposes stable model availability, not provider health.
// Explicit configuration remains available through transient provider failures.
func CapabilitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "GET required", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	json.NewEncoder(w).Encode(map[string]bool{"model": model.Available()})
}

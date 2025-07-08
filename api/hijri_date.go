package api

import (
	"encoding/json"
	"net/http"

	"github.com/asim/reminder/hijri"
)

func HijriDateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hijri.Date())
}

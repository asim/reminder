package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	hijri "github.com/hablullah/go-hijri"
)

type HijriDateResponse struct {
	Date   string `json:"date"`
	Display string `json:"display"`
}

func HijriDateHandler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	h, err := hijri.CreateUmmAlQuraDate(now)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to calculate Hijri date"}`))
		return
	}

	dateStr := fmt.Sprintf("%02d-%02d-%04d", h.Day, h.Month, h.Year)

	ordinal := func(n int) string {
		if n == 1 {
			return "st"
		} else if n == 2 {
			return "nd"
		} else if n == 3 {
			return "rd"
		} else if n%10 == 1 && n%100 != 11 {
			return "st"
		} else if n%10 == 2 && n%100 != 12 {
			return "nd"
		} else if n%10 == 3 && n%100 != 13 {
			return "rd"
		}
		return "th"
	}
	months := []string{"Muharram", "Safar", "Rabiʿ al-awwal", "Rabiʿ al-thani", "Jumada al-awwal", "Jumada al-thani", "Rajab", "Shaʿban", "Ramadan", "Shawwal", "Dhu al-Qiʿdah", "Dhu al-Hijjah"}
	display := "Today is the " +
		fmt.Sprintf("%d", h.Day) +
		ordinal(int(h.Day)) +
		" of " + months[int(h.Month)-1] + ", " + fmt.Sprintf("%d", h.Year)

	resp := HijriDateResponse{
		Date:   dateStr,
		Display: display,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

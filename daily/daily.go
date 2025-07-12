package daily

import (
	"fmt"
	"time"

	hijri "github.com/hablullah/go-hijri"
)

type Today struct {
	Date    string `json:"date"`
	Hijri   string `json:"hijri"`
	Display string `json:"display"`
}

func Date() *Today {
	now := time.Now()
	h, err := hijri.CreateUmmAlQuraDate(now)
	if err != nil {
		return new(Today)
	}
	ordinal := func(n int64) string {
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

	dateStr := fmt.Sprintf("%04d-%02d-%02d", h.Year, h.Month, h.Day)

	display := fmt.Sprintf("%d", h.Day) +
		ordinal(h.Day) +
		" of " + months[int(h.Month)-1] + ", " + fmt.Sprintf("%d", h.Year)

	return &Today{
		Date:    now.Format("2006-01-02"),
		Hijri:   dateStr,
		Display: display,
	}
}

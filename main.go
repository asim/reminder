package main

import (
	"fmt"
	"os"

	"github.com/asim/reminder/quran"
)

var md = func(q map[string]*quran.Chapter) string {
	var data string

	// 114 surahs
	for i := 0; i < 114; i++ {
		ch := q[fmt.Sprintf("%d", i+1)]

		data += fmt.Sprintln()
		data += fmt.Sprintf(`# %d`, ch.Number)
		data += fmt.Sprintln()
		data += fmt.Sprintln()
		data += fmt.Sprintf(`### %s`, ch.Name)
		data += fmt.Sprintln()

		// max 286 ayahs
		for _, verse := range ch.Verses {
			data += fmt.Sprintln()
			data += fmt.Sprintf(verse.Text)
			data += fmt.Sprintln()
		}

	}

	return data
}

func main() {
	fmt.Println("Loading source")
	q := quran.Load()

	fmt.Println("Compiling text")
	text := md(q)

	fmt.Println("Saving text")
	os.WriteFile("index.md", []byte(text), 0644)
}

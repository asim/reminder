package main

import (
	"fmt"
	"os"

	"github.com/asim/reminder/quran"
)

func main() {
	fmt.Println("Loading text")
	text := quran.Markdown()

	fmt.Println("Saving text")
	os.WriteFile("index.md", []byte(text), 0644)
}

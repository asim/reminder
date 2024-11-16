package files

import (
	"embed"
)

//go:embed *.html
var files embed.FS

func Get(name string) string {
	f, err := files.ReadFile(name + ".html")
	if err != nil {
		return ""
	}

	return string(f)
}

package files

import (
	"embed"
)

//go:embed *.html
//go:embed *.json
//go:embed *.otf
var files embed.FS

func Get(name string) string {
	f, err := files.ReadFile(name)
	if err != nil {
		return ""
	}

	return string(f)
}

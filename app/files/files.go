package files

import (
	"embed"
)

//go:embed *.html
//go:embed *.otf
//go:embed *.js
//go:embed *.png
//go:embed manifest.webmanifest
var files embed.FS

func Get(name string) string {
	f, err := files.ReadFile(name)
	if err != nil {
		return ""
	}

	return string(f)
}

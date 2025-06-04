package app

import (
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

// FSHandler404 provides the function signature for passing to the FileServerWith404
type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

/*
FileServerWith404 wraps the http.FileServer checking to see if the url path exists first.
If the file fails to exist it calls the supplied FSHandle404 function
The implementation can choose to either modify the request, e.g. change the URL path and return true to have the
default FileServer handling to still take place, or return false to stop further processing, for example if you wanted
to write a custom response
e.g. redirects to root and continues the file serving handler chain

	func fileSystem404(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool) {
		//if not found redirect to /
		r.URL.Path = "/"
		return true
	}

Use the same as you would with a http.FileServer e.g.

	r.Handle("/", http.StripPrefix("/", mw.FileServerWith404(http.Dir("./staticDir"), fileSystem404)))
*/
func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure the url path starts with /
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				// If path ends with a trailing slash, try again without it
				if len(upath) > 1 && strings.HasSuffix(upath, "/") {
					altPath := strings.TrimSuffix(upath, "/")
					f2, err2 := root.Open(altPath)
					if err2 == nil {
						f2.Close()
						// Rewrite the request path and serve
						r2 := new(http.Request)
						*r2 = *r
						urlCopy := *r.URL
						urlCopy.Path = altPath
						r2.URL = &urlCopy
						fs.ServeHTTP(w, r2)
						return
					}
				}
				log.Printf("Embedded asset not found: %s", upath)
				// call handler
				if handler404 != nil {
					doDefault := handler404(w, r)
					if !doDefault {
						return
					}
				}
			} else {
				log.Printf("error opening asset %s: %v", upath, err)
			}
		}

		// close if successfully opened
		if err == nil {
			f.Close()
		}

		// default serve
		fs.ServeHTTP(w, r)
	})
}

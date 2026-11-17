// Package corpusversion fingerprints embedded source data independently of builds.
package corpusversion

import (
	"crypto/sha256"
	"fmt"
	"io/fs"
)

func Sum(files fs.FS) string {
	h := sha256.New()
	err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}
		fmt.Fprintf(h, "%d:%s:%d:", len(path), path, len(data))
		h.Write(data)
		return nil
	})
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

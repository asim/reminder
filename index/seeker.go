package index

import (
	"errors"
	"io"
	"io/fs"
)

type ReadSeekerWrapper struct {
	data   []byte
	offset int64
}

// NewReadSeekerWrapper creates a ReadSeeker from an fs.File
func NewReadSeekerWrapper(f fs.File) (io.ReadSeeker, error) {
	data, err := io.ReadAll(f) // Read the entire content
	if err != nil {
		return nil, err
	}
	return &ReadSeekerWrapper{data: data}, nil
}

// Read reads data into p
func (r *ReadSeekerWrapper) Read(p []byte) (int, error) {
	if r.offset >= int64(len(r.data)) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.offset:])
	r.offset += int64(n)
	return n, nil
}

// Seek sets the offset for the next Read
func (r *ReadSeekerWrapper) Seek(offset int64, whence int) (int64, error) {
	var newOffset int64
	switch whence {
	case io.SeekStart:
		newOffset = offset
	case io.SeekCurrent:
		newOffset = r.offset + offset
	case io.SeekEnd:
		newOffset = int64(len(r.data)) + offset
	default:
		return 0, errors.New("invalid whence")
	}
	if newOffset < 0 || newOffset > int64(len(r.data)) {
		return 0, errors.New("invalid offset")
	}
	r.offset = newOffset
	return r.offset, nil
}

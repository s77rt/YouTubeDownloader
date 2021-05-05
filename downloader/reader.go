package downloader

import (
	"io"
)

type ReaderWithCounter struct {
	io.Reader
	Transferred int64
}

func (r *ReaderWithCounter) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.Transferred += int64(n)

	return n, err
}

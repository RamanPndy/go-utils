package goutils

import "io"

type bytesReader string

func (b bytesReader) Read(p []byte) (int, error) {
	n := copy(p, []byte(string(b)))
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func NewBytesReader(s string) io.Reader {
	return bytesReader(s)
}

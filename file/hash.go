package file

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

var HashType = md5.New()

type HashWriter struct {
	io.Writer
	hash hash.Hash
}

func NewHashWrite(writer io.Writer) *HashWriter {
	return &HashWriter{
		Writer: writer,
		hash:   HashType,
	}
}

func (w *HashWriter) Write(p []byte) (n int, err error) {
	n, err = w.Writer.Write(p)
	if n > 0 {
		w.hash.Write(p)
	}
	return n, err
}

func (w *HashWriter) Hash() string {
	b := w.hash.Sum(nil)
	return hex.EncodeToString(b)
}

type HashReader struct {
	io.Reader
	hash hash.Hash
}

func NewHashReader(reader io.Reader) *HashReader {
	return &HashReader{
		Reader: reader,
		hash:   HashType,
	}
}

func (r *HashReader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	if n > 0 {
		r.hash.Write(p)
	}
	return n, err
}

func (r *HashReader) Hash() string {
	b := r.hash.Sum(nil)
	return hex.EncodeToString(b)
}

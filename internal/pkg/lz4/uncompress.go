package lz4

import (
	"io"
	"os"
	"strings"

	"github.com/pierrec/lz4/v4"
)

func Uncompress(zfilePath string) (string, error) {
	zr := lz4.NewReader(nil)

	zfile, err := os.Open(zfilePath)
	if err != nil {
		return "", err
	}
	zinfo, err := zfile.Stat()
	if err != nil {
		return "", err
	}
	mode := zinfo.Mode() // use the same mode for the output file

	filename := strings.TrimSuffix(zfilePath, lz4Extension)
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, mode)
	if err != nil {
		return "", err
	}
	zr.Reset(zfile)

	var out io.Writer = file

	_, err = io.Copy(out, zr)
	if err != nil {
		return "", err
	}
	for _, c := range []io.Closer{zfile, file} {
		err := c.Close()
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}

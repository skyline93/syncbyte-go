package storageunit

import (
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	Path string
}

func (l *Local) Save(fileName string, read io.Reader) error {
	p := filepath.Join(l.Path, fileName)

	fs, err := os.Create(p)
	if err != nil {
		return err
	}
	defer fs.Close()

	_, err = io.Copy(fs, read)
	if err != nil {
		return err
	}

	return nil
}

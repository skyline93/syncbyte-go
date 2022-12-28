package file

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4"
)

var blockSize int64 = 1024 * 1024

func UploadBigFile(filename string, callback func(string, []byte) error) error {
	var (
		buf        []byte = make([]byte, blockSize)
		compressed []byte = make([]byte, blockSize*2)
		partNum    int    = 1
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for {
		hReader := NewHashReader(file)
		rSize, err := hReader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		l, err := lz4.CompressBlock(buf[:rSize], compressed, nil)
		if err != nil {
			return err
		}

		if err := callback(hReader.Hash(), compressed[:l]); err != nil {
			return err
		}

		partNum += 1
	}

	return nil
}

func UploadSmallFile(filename string, callback func(string, []byte) error) error {
	var (
		buf        []byte = make([]byte, blockSize)
		compressed []byte = make([]byte, blockSize*2)
	)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	hReader := NewHashReader(file)

	rSize, err := hReader.Read(buf)
	if err != nil {
		if err != io.EOF {
			return err
		}
	}

	l, err := lz4.CompressBlock(buf[:rSize], compressed, nil)
	if err != nil {
		return err
	}

	if err := callback(hReader.Hash(), compressed[:l]); err != nil {
		return err
	}

	return nil
}

type DataStor interface {
	UploadFile(filename string) error
}

type NasVolume struct {
	MountPoint string
}

func NewNasVolume(mountPoint string) *NasVolume {
	return &NasVolume{MountPoint: mountPoint}
}

func (n *NasVolume) UploadFile(filename string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if fi.Size() >= blockSize {
		return UploadBigFile(filename, n.uploadFile)
	}

	return UploadSmallFile(filename, n.uploadFile)
}

func (n *NasVolume) uploadFile(filename string, buf []byte) error {
	file, err := os.Create(filepath.Join(n.MountPoint, filename))
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, bytes.NewReader(buf))
	if err != nil {
		return err
	}

	return nil
}

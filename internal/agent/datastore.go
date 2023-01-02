package agent

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/pierrec/lz4"
)

var blockSize int64 = 1024 * 1024

type PartInfo struct {
	Index int    `json:"index"`
	MD5   string `json:"md5"`
	Size  int64  `json:"size"`
}

type FileInfo struct {
	Name      string      `json:"name"`
	Size      int64       `json:"size"`
	MD5       string      `json:"md5"`
	PartInfos []*PartInfo `json:"part_info"`
}

func (fi *FileInfo) String() string {
	v, _ := json.Marshal(fi)
	return string(v)
}

func UploadBigFile(filename string, callback func(string, []byte) error, fileInfo *FileInfo) error {
	var (
		buf        []byte = make([]byte, blockSize)
		compressed []byte = make([]byte, blockSize*2)
		partIndex  int    = 1
		partInfos  []*PartInfo
	)

	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fhReader := NewHashReader(file)

	for {
		hReader := NewHashReader(fhReader)
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

		pi := &PartInfo{Index: partIndex, MD5: hReader.Hash(), Size: int64(rSize)}
		partInfos = append(partInfos, pi)

		partIndex += 1
	}

	fileInfo.Name = fi.Name()
	fileInfo.Size = fi.Size()
	fileInfo.MD5 = fhReader.Hash()
	fileInfo.PartInfos = partInfos

	return nil
}

func UploadSmallFile(filename string, callback func(string, []byte) error, fileInfo *FileInfo) error {
	var (
		buf        []byte = make([]byte, blockSize)
		compressed []byte = make([]byte, blockSize*2)
	)

	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

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

	fileInfo.Name = fi.Name()
	fileInfo.Size = fi.Size()
	fileInfo.MD5 = hReader.Hash()

	return nil
}

type DataStore interface {
	UploadFile(filename string, fileInfo *FileInfo) error
}

type NasVolume struct {
	MountPoint string
}

func NewNasVolume(mountPoint string) *NasVolume {
	return &NasVolume{MountPoint: mountPoint}
}

func (n *NasVolume) UploadFile(filename string, fileInfo *FileInfo) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}

	if fi.Size() >= blockSize {
		return UploadBigFile(filename, n.uploadFile, fileInfo)
	}

	return UploadSmallFile(filename, n.uploadFile, fileInfo)
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

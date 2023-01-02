package agent

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/pierrec/lz4"
)

var blockSize int64 = 1024 * 1024

type PartInfo struct {
	Index int    `json:"index"`
	MD5   string `json:"md5"`
	Size  int64  `json:"size"`
}

type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	GID        uint32 `json:"gid"`
	UID        uint32 `json:"uid"`
	Device     uint64 `json:"device"`
	DeviceID   uint64 `json:"device_id"`
	BlockSize  int64  `json:"block_size"`
	Blocks     int64  `json:"blocks"`
	AccessTime int64  `json:"atime"`
	ModTime    int64  `json:"mtime"`
	ChangeTime int64  `json:"ctime"`

	Path      string      `json:"path"`
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
		rSize, err := fhReader.Read(buf)
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

		h := md5.New()
		h.Write(buf)
		hash := hex.EncodeToString(h.Sum(nil))

		if err := callback(hash, compressed[:l]); err != nil {
			return err
		}

		pi := &PartInfo{Index: partIndex, MD5: hash, Size: int64(rSize)}
		partInfos = append(partInfos, pi)

		partIndex += 1
	}

	s := fi.Sys().(*syscall.Stat_t)

	fileInfo.Name = fi.Name()
	fileInfo.Size = fi.Size()
	fileInfo.GID = s.Uid
	fileInfo.UID = s.Uid
	fileInfo.Device = s.Rdev
	fileInfo.DeviceID = s.Dev
	fileInfo.BlockSize = int64(s.Blksize)
	fileInfo.Blocks = s.Blocks
	fileInfo.AccessTime = s.Atim.Nano()
	fileInfo.ModTime = s.Mtim.Nano()
	fileInfo.ChangeTime = s.Ctim.Nano()

	fileInfo.Path = filename
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

	s := fi.Sys().(*syscall.Stat_t)

	fileInfo.Name = fi.Name()
	fileInfo.Size = fi.Size()
	fileInfo.GID = s.Uid
	fileInfo.UID = s.Uid
	fileInfo.Device = s.Rdev
	fileInfo.DeviceID = s.Dev
	fileInfo.BlockSize = int64(s.Blksize)
	fileInfo.Blocks = s.Blocks
	fileInfo.AccessTime = s.Atim.Nano()
	fileInfo.ModTime = s.Mtim.Nano()
	fileInfo.ChangeTime = s.Ctim.Nano()

	fileInfo.Path = filename
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

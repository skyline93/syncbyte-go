package backuper

import (
	"archive/tar"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/skyline93/syncbyte-go/internal/pkg/lz4"
)

func tarFile(basepath, filePath string, fi fs.FileInfo, tw *tar.Writer, r io.Reader) error {
	header := new(tar.Header)

	path, err := filepath.Rel(filepath.Join(basepath, "../"), filePath)
	if err != nil {
		return err
	}

	header.Name = path
	header.Size = fi.Size()
	header.Mode = int64(fi.Mode())
	header.ModTime = fi.ModTime()

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tw, r); err != nil {
		return err
	}

	return nil
}

func tarDir(dir string, dest io.WriteCloser) error {
	fChan := make(chan string)
	go scanDir(dir, fChan)

	tw := tar.NewWriter(dest)
	defer tw.Close()

	func() {
		for f := range fChan {
			fi, err := os.Stat(f)
			if err != nil {
				log.Printf("stat file [%s] error, msg: %v", f, err)
				continue
			}

			fp, err := os.Open(f)
			if err != nil {
				log.Printf("open file [%s] failed", f)
				continue
			}

			if err := tarFile(dir, f, fi, tw, fp); err != nil {
				log.Printf("tar file [%s] failed, msg: %v", f, err)
				fp.Close()
				continue
			}
			fp.Close()
		}
	}()

	return nil
}

func removePrefixPath(path string) string {
	v := strings.Split(path, "/")

	if v[0] == "" {
		v = v[2:]
	} else {
		v = v[1:]
	}

	path = strings.Join(v, "/")
	return path
}

func untarFile(reader io.Reader, dst string) error {
	tr := tar.NewReader(reader)

	dst, err := filepath.Abs(dst)
	if err != nil {
		return err
	}

	for {
		header, err := tr.Next()
		switch {
		// no more files
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue
		}

		target := filepath.Join(dst, removePrefixPath(header.Name))
		dir := filepath.Dir(target)

		_, err = os.Stat(dir)
		if os.IsNotExist(err) {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		switch header.Typeflag {
		// create directory if doesn't exit
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			fmt.Printf("open file %s\n", target)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				log.Printf("open file faile, target: %s", target)
				return err
			}

			// copy contents to file
			if _, err := io.Copy(f, tr); err != nil {
				f.Close()
				return err
			}

			f.Close()
		}
	}
}

func Compress(dir string, destName string) (string, error) {
	path := filepath.Join(filepath.Join(dir, "../"), destName+".tar")

	df, err := os.Create(path)
	if err != nil {
		return "", err
	}

	if err := tarDir(dir, df); err != nil {
		return "", err
	}

	dest, err := lz4.Compress(path)
	if err != nil {
		return "", err
	}

	if err := os.Remove(path); err != nil {
		return "", err
	}

	return dest, nil
}

func Uncompress(filePath string, destPath string) error {
	destFile, err := lz4.Uncompress(filePath)
	if err != nil {
		return err
	}

	fp, err := os.Open(destFile)
	if err != nil {
		return err
	}

	if err := untarFile(fp, destPath); err != nil {
		return err
	}

	fp.Close()

	if err := os.Remove(destFile); err != nil {
		return err
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}

	return nil
}

func scanDir(root string, fChan chan string) error {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fChan <- path
		return nil
	})

	close(fChan)
	return err
}

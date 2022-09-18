package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func Compress(file string, isDel bool) (string, error) {
	gzipFile := fmt.Sprintf("%s.gz", file)

	fr, err := os.Open(file)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadAll(fr)
	if err != nil {
		return "", err
	}

	defer fr.Close()

	outputFile, err := os.Create(gzipFile)
	if err != nil {
		return "", err
	}

	gzipWriter := gzip.NewWriter(outputFile)

	_, err = gzipWriter.Write(data)
	if err != nil {
		return "", err
	}

	defer gzipWriter.Close()

	if isDel {
		if err = os.Remove(file); err != nil {
			return "", err
		}
	}
	return gzipFile, err
}

func UnCompress(file string, isDel bool) (string, error) {
	unzipFile := strings.Trim(file, ".gz")

	gzipFile, err := os.Open(file)
	if err != nil {
		return "", err
	}

	gzipReader, err := gzip.NewReader(gzipFile)
	if err != nil {
		return "", err
	}
	defer gzipReader.Close()

	outfileWriter, err := os.Create(unzipFile)
	if err != nil {
		return "", err
	}
	defer outfileWriter.Close()

	_, err = io.Copy(outfileWriter, gzipReader)
	if err != nil {
		return "", err
	}

	if isDel {
		if err = os.Remove(file); err != nil {
			return "", err
		}
	}
	return unzipFile, err
}

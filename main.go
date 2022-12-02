package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func main() {
	// dir := "/home/greene/workspace/syncbyte-go/output"
	file := "/home/greene/workspace/syncbyte-go/output/testdir1/syncbyte-agent.txt"

	fn := path.Base(file)
	fmt.Println(fn)

	// fChan := make(chan string)

	// go func() {
	// 	for f := range fChan {
	// 		log.Printf("open file [%s]", f)

	// 	}
	// }()

	// if err := ScanDir(dir, fChan); err != nil {
	// 	panic(err)
	// }
}

func ScanDir(root string, fChan chan string) error {
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

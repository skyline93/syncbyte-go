package main

import "github.com/skyline93/syncbyte-go/internal/agent/backuper"

func main() {
	// dir := "/home/greene/workspace/syncbyte-go/output.bak"

	// v := filepath.Join("workspace/syncbyte-go/output", "../")
	// // v := filepath.Join(dir, "../")
	// fmt.Printf("v: %s\n", v)

	// v2 := strings.Split("workspace/syncbyte-go/output", "/")
	// fmt.Printf("v2: %v\n", v2)

	// dir = "output"
	// v3 := strings.Split(dir, "/")
	// fmt.Printf("v3: %v\n", v3)

	// if v3[0] == "" {
	// 	v3 = v3[2:]
	// } else {
	// 	v3 = v3[1:]
	// }
	// v4 := strings.Join(v3, "/")

	// fmt.Printf("v4: %v\n", v4)

	// basepath := filepath.Join(dir, "../")

	// fmt.Println(filepath.Join(dir, "../"))

	// path, err := filepath.Rel(basepath, "/home/greene/workspace/syncbyte-go/output/syncbyte")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("path: %s", path)

	// dest := "/home/greene/workspace/syncbyte-go/datastore"

	// file := "/home/greene/workspace/syncbyte-go/output/syncbyte.txt"
	// file := "/home/greene/workspace/syncbyte-go/Makefile"

	// fi, err := os.Stat(file)
	// if err != nil {
	// 	panic(err)
	// }

	// sf, err := os.Open(file)
	// if err != nil {
	// 	panic(err)
	// }

	// defer sf.Close()

	//==
	// df, err := os.Create("abc.tar")
	// if err != nil {
	// 	panic(err)
	// }

	// defer df.Close()

	// if err := backuper.TarDir(dir, df); err != nil {
	// 	panic(err)
	// }

	// // df.Close()

	// if err := lz4.Compress("abc.tar"); err != nil {
	// 	panic(err)
	// }

	// if err := lz4.Uncompress("abc.tar.lz4"); err != nil {
	// 	panic(err)
	// }

	// cf, err := os.Create("")
	// backuper.Compress()

	// =========

	// if err := backuper.Compress(dir, "12345"); err != nil {
	// 	panic(err)
	// }

	// v := filepath.Join(filepath.Join(dir, "../"), "12345.tar.lz4")
	// if err := backuper.Uncompress(v, "1234"); err != nil {
	// 	panic(err)
	// }

	// =========

	// if err := TarFile(file, fi, df, sf); err != nil {
	// 	panic(err)
	// }

	dir := "/home/greene/workspace/syncbyte-go/internal"
	dest := "/home/greene/workspace/syncbyte-go/datastore"

	dataStore := backuper.NewLocalFileSystem(dest)
	b := backuper.NewBackuper(true, dataStore)

	b.Backup(dir)
}

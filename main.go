package main

// Connection URI
// const uri = "mongodb://mongoadmin:123456@10.168.1.202:27017/?maxPoolSize=20&w=majority"

type FileInfo struct {
	Name string
	Size int64
}

func main() {
	// client, err := mongodb.NewClient(uri)
	// if err != nil {
	// 	panic(err)
	// }
	// defer client.Close()

	// col := client.GetCollection("backup01")

	// col.InsertOne(context.TODO())
}

package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error while loading environment files")
	}
}

func main() {
	// create Bucket
	// create Entry
	// upload Content
	// delete Bucket
	clientProvider := ClientProvider{}
	contentClient, err := clientProvider.GetClient()
	if err != nil {
		log.Fatalln(err)
	}
	_, _ = contentClient.CreateBucket(os.Getenv("PROJECT_ID"))
}

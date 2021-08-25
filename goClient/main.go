package main

import (
	"contentWorkflow/utils"
	"log"

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
	bucketId, err := contentClient.CreateBucket()
	if err != nil {
		log.Fatalln(err)
	}

	content, err := utils.GenerateRandomContent(100)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("RANDOM GENERATED CONTENT\n%v\n\n", content)

	entryId, err := contentClient.CreateEntry(bucketId, content)
	if err != nil {
		log.Fatalln(err)
	}

	contentClient.UploadContent(bucketId, entryId, content)

	contentClient.CreateRelease(bucketId)
	// contentClient.DeleteBucket(bucketId)
}

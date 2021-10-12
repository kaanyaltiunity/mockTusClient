package main

import "contentWorkflow/utils"

type ContentClient interface {
	CreateBucket() (string, error)
	CreateEntry(bucketId string, content *utils.Content) (string, string, error)
	UploadContent(bucketId string, entryId string, content *utils.Content)
	CreateRelease(bucketId string) string
	DeleteBucket(bucketId string)
}

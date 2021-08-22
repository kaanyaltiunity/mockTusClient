package main

import "contentWorkflow/utils"

type ContentClient interface {
	CreateBucket() (string, error)
	CreateEntry(bucketId string, content *utils.Content) (string, error)
	UploadContent(bucketId string, entryId string) error
	DeleteBucket(bucketId string) error
}

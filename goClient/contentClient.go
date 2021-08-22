package main

type ContentClient interface {
	CreateBucket(projectId string) (string, error)
	CreateEntry(bucketId string) (string, error)
	UploadContent(bucketId string) error
	DeleteBucket(bucketId string) error
}

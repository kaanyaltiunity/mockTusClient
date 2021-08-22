package gatewayClient

import (
	"bytes"
	"contentWorkflow/baseClient"
	"contentWorkflow/payloads"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

const (
	stagingBaseUrl string = "https://staging.services.unity.com/api/ccd/management/v1/"
)

type GatewayClient struct {
	baseUrl string
	baseClient.BaseClient
}

func NewGatewayClient() *GatewayClient {
	log.Println("CREATING GATEWAY CLIENT\n")
	base := baseClient.NewBaseClient()
	client := GatewayClient{
		BaseClient: base,
		baseUrl:    stagingBaseUrl,
	}
	client.SetAuth("Bearer", os.Getenv("BEARER_TOKEN"))
	return &client
}

func (gatewayClient *GatewayClient) CreateBucket(projectId string) (string, error) {
	uuid := uuid.NewString()
	payload := payloads.CreateBucket{
		Description: "testDescription",
		Name:        fmt.Sprintf("testBucket-%s", uuid),
		ProjectGuid: projectId,
	}

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("GATEWAY BUCKET CREATION MARSHALLED PAYLOAD %s\n\n", string(marshalledPayload))

	request, err := http.NewRequest("POST", fmt.Sprintf("%sprojects/%s/buckets", stagingBaseUrl, projectId), bytes.NewBuffer(marshalledPayload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", gatewayClient.Auth)

	log.Printf("GATEWAY CREATE BUCKET HEADERS %v\n\n", request.Header)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("BUCKET CREATION RESPONSE %v\n\n", response)

	return "test", nil
}
func (gatewayClient *GatewayClient) CreateEntry(bucketId string) (string, error) {
	return "", nil
}
func (gatewayClient *GatewayClient) UploadContent(bucketId string) error {
	return nil
}
func (gatewayClient *GatewayClient) DeleteBucket(bucketId string) error {
	return nil
}

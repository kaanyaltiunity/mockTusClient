package gatewayClient

import (
	"bytes"
	"contentWorkflow/baseClient"
	"contentWorkflow/payloads"
	"contentWorkflow/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
)

var baseUrls map[string]string = map[string]string{
	"staging": "https://staging.services.unity.com/api/ccd/management/v1/",
	"dev":     "localhost:9000/api/ccd/management/v1/",
}

type GatewayClient struct {
	baseUrl   string
	projectId string
	baseClient.BaseClient
}

func NewGatewayClient(projectId string) *GatewayClient {
	log.Println("CREATING GATEWAY CLIENT\n")
	base := baseClient.NewBaseClient()
	client := GatewayClient{
		BaseClient: base,
		baseUrl:    baseUrls[os.Getenv("ENV")],
		projectId:  projectId,
	}
	client.SetAuth("Bearer", os.Getenv("BEARER_TOKEN"))
	return &client
}

func (gatewayClient *GatewayClient) CreateBucket() (string, error) {
	uuid := uuid.NewString()
	payload := payloads.CreateBucket{
		Description: "testDescription",
		Name:        fmt.Sprintf("testBucket-%s", uuid),
		ProjectGuid: gatewayClient.projectId,
	}

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("GATEWAY BUCKET CREATION MARSHALLED PAYLOAD\n%s\n\n", string(marshalledPayload))

	request, err := http.NewRequest("POST", fmt.Sprintf("%sprojects/%s/buckets", gatewayClient.baseUrl, gatewayClient.projectId), bytes.NewBuffer(marshalledPayload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", gatewayClient.Auth)

	log.Printf("GATEWAY CREATE BUCKET HEADERS\n%v\n\n", request.Header)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("BUCKET CREATION RESPONSE STATUS\n%v\n\n", response.Status)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var unmarshalledResponseBody map[string]interface{}

	err = json.Unmarshal(responseBody, &unmarshalledResponseBody)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("BUCKET CREATION RESPONSE BODY\n%v\n\n", unmarshalledResponseBody)
	log.Printf("BUCKET CREATION BUCKET ID\n%v\n\n", unmarshalledResponseBody["id"])

	return unmarshalledResponseBody["id"].(string), nil
}
func (gatewayClient *GatewayClient) CreateEntry(bucketId string, content *utils.Content) (string, error) {
	contentPath := fmt.Sprintf("test_entry_%s", utils.RandomString(10))
	log.Printf("CONTENT PATH\n%s\n\n", contentPath)

	payload := payloads.CreateEntry{
		Path:        contentPath,
		ContentHash: content.Hash,
		ContentSize: content.Size,
		ContentType: content.Type,
	}

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("ENTRY CREATION MARSHALLED PAYLOAD\n%s\n\n", string(marshalledPayload))

	request, err := http.NewRequest("POST", fmt.Sprintf("%sprojects/%s/buckets/%s/entries", gatewayClient.baseUrl, gatewayClient.projectId, bucketId), bytes.NewBuffer(marshalledPayload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", gatewayClient.Auth)

	log.Printf("GATEWAY ENTRY CREATION HEADERS\n%v\n\n", request.Header)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("ENTRY CREATION RESPONSE STATUS\n%v\n\n", response.Status)
	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var unmarshalledResponseBody map[string]interface{}

	err = json.Unmarshal(responseBody, &unmarshalledResponseBody)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("ENTRY CREATION RESPONSE BODY\n%v\n\n", unmarshalledResponseBody)
	log.Printf("ENTRY CREATION ENTRY ID\n%v\n\n", unmarshalledResponseBody["entryid"])

	return unmarshalledResponseBody["entryid"].(string), nil
}
func (gatewayClient *GatewayClient) UploadContent(bucketId string, entryId string, content *utils.Content) {
	request, err := http.NewRequest("PATCH", fmt.Sprintf("%sprojects/%s/buckets/%s/entries/%s/content", gatewayClient.baseUrl, gatewayClient.projectId, bucketId, entryId), bytes.NewBuffer(content.Bytes))
	request.Header.Set("Content-Type", "application/offset+octet-stream")
	request.Header.Set("Content-Length", strconv.FormatInt(int64(content.Size), 10))
	request.Header.Set("Authorization", gatewayClient.Auth)
	request.Header.Set("Upload-Offset", "0")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("GATEWAY CONTENT UPLOAD HEADERS\n%v\n\n", request.Header)

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("CONTENT UPLOAD RESPONSE STATUS\n%v\n\n", response.Status)
}

func (gatewayClient *GatewayClient) CreateRelease(bucketId string) string {
	body, _ := json.Marshal(map[string]interface{}{})
	request, err := http.NewRequest("POST", fmt.Sprintf("%sprojects/%s/buckets/%s/releases", gatewayClient.baseUrl, gatewayClient.projectId, bucketId), bytes.NewBuffer(body))
	request.Header.Set("Authorization", gatewayClient.Auth)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("GATEWAY CREATE RELEASE URL\n%v\n\n", request.URL)
	log.Printf("GATEWAY CREATE RELEASE HEADERS\n%v\n\n", request.Header)

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("CREATE RELEASE RESPONSE STATUS\n%v\n\n", response.Status)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var unmarshalledResponseBody map[string]interface{}

	err = json.Unmarshal(responseBody, &unmarshalledResponseBody)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("CREATE RELEASE RESPONSE BODY\n%v\n\n", unmarshalledResponseBody)
	return unmarshalledResponseBody["releaseid"].(string)
}

func (gatewayClient *GatewayClient) DeleteBucket(bucketId string) {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%sprojects/%s/buckets/%s", gatewayClient.baseUrl, gatewayClient.projectId, bucketId), nil)
	request.Header.Set("Authorization", gatewayClient.Auth)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := gatewayClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("BUCKET DELETE RESPONSE STATUS\n%v\n\n", response.Status)
}

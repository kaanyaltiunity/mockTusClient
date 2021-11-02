package cdsClient

import (
	"bytes"
	"contentWorkflow/baseClient"
	"contentWorkflow/payloads"
	"contentWorkflow/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/eventials/go-tus"
	"github.com/google/uuid"
)

var baseUrls map[string]string = map[string]string{
	"staging": "https://content-api-stg.cloud.unity3d.com/api/v1/",
	"dev":     "http://localhost:30332/api/v1/",
}

type CdsClient struct {
	baseUrl   string
	projectId string
	baseClient.BaseClient
}

func NewCdsClient(projectId string) *CdsClient {
	log.Println("CREATING CDS CLIENT\n")
	base := baseClient.NewBaseClient()
	env := os.Getenv("ENV")
	client := CdsClient{
		BaseClient: base,
		baseUrl:    baseUrls[env],
		projectId:  projectId,
	}
	apiKeyVar := fmt.Sprintf("API_KEY_%s", strings.ToUpper(env))
	apiKey := os.Getenv(apiKeyVar)
	fmt.Printf("CDS CLIENT API KEY: %s\n", apiKey)
	client.SetAuth("Basic", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf(":%s", apiKey))))
	return &client
}

func (cdsClient *CdsClient) CreateBucket() (string, error) {
	uuid := uuid.NewString()
	payload := payloads.CreateBucket{
		Description: "testDescription",
		Name:        fmt.Sprintf("testBucket-%s", uuid),
		ProjectGuid: cdsClient.projectId,
	}

	marshalledPayload, err := json.Marshal(payload)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("CDS BUCKET CREATION MARSHALLED PAYLOAD\n%s\n\n", string(marshalledPayload))

	request, err := http.NewRequest("POST", fmt.Sprintf("%sprojects/%s/buckets/", cdsClient.baseUrl, cdsClient.projectId), bytes.NewBuffer(marshalledPayload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", cdsClient.Auth)

	log.Printf("CDS CREATE BUCKET HEADERS\n%v\n\n", request.Header)
	log.Printf("CDS CREATE BUCKET URL\n%v\n\n", request.URL)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := cdsClient.Do(request)
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

func (cdsClient *CdsClient) CreateEntry(bucketId string, content *utils.Content) (string, string, error) {
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

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/buckets/%s/entries/", cdsClient.baseUrl, bucketId), bytes.NewBuffer(marshalledPayload))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", cdsClient.Auth)

	log.Printf("CDS ENTRY CREATION HEADERS\n%v\n\n", request.Header)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := cdsClient.Do(request)
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

	return unmarshalledResponseBody["entryid"].(string), unmarshalledResponseBody["current_versionid"].(string), nil
}

func (cdsClient *CdsClient) UploadContent(bucketId string, entryId string, content *utils.Content) {
	path := fmt.Sprintf("%s/buckets/%s/entries/%s/content/", cdsClient.baseUrl, bucketId, entryId)

	goTusEnabled, err := strconv.ParseBool(os.Getenv("GO_TUS_ENABLED"))
	if err != nil {
		log.Fatalln(err)
	}

	if goTusEnabled {
		cdsClient.uploadWithGoTus(path, content)
	} else {
		cdsClient.uploadWithHttp(path, content)
	}
}

func (cdsClient *CdsClient) uploadWithHttp(path string, content *utils.Content) {
	request, err := http.NewRequest("PATCH", path, bytes.NewBuffer(content.Bytes))
	request.Header.Set("Content-Type", "application/offset+octet-stream")
	request.Header.Set("Content-Length", strconv.FormatInt(int64(content.Size), 10))
	request.Header.Set("Authorization", cdsClient.Auth)
	request.Header.Set("Upload-Offset", "0")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("CDS CONTENT UPLOAD HEADERS\n%v\n\n", request.Header)

	response, err := cdsClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("CONTENT UPLOAD RESPONSE STATUS\n%v\n\n", response.Status)
}

func (cdsClient *CdsClient) uploadWithGoTus(path string, content *utils.Content) {
	fmt.Println("UPLOADING WITH GO TUS\n")

	tusConfig := tus.Config{
		Header: http.Header{
			"Authorization": []string{cdsClient.Auth},
		},
		ChunkSize: 5,
	}

	client, err := tus.NewClient(path, &tusConfig)
	if err != nil {
		log.Fatalln(err)
	}

	upload := tus.NewUploadFromBytes(content.Bytes)

	uploader, err := client.CreateUpload(upload)
	if err != nil {
		log.Fatalln(err)
	}

	err = uploader.Upload()
	if err != nil {
		log.Fatalln(err)
	}
}

func (cdsClient *CdsClient) CreateRelease(bucketId string) string {
	body, _ := json.Marshal(map[string]interface{}{})
	request, err := http.NewRequest("POST", fmt.Sprintf("%sbuckets/%s/releases/", cdsClient.baseUrl, bucketId), bytes.NewBuffer(body))
	request.Header.Set("Authorization", cdsClient.Auth)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("CDS CREATE RELEASE HEADERS\n%v\n\n", request.Header)
	log.Printf("CDS CREATE RELEASE URL\n%v\n\n", request.URL)
	response, err := cdsClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("CREATE RELEASE RESPONSE STATUS\n%v\n\n", response.Status)
	return ""
}

func (cdsClient *CdsClient) DeleteBucket(bucketId string) {
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/buckets/%s/", cdsClient.baseUrl, bucketId), nil)
	request.Header.Set("Authorization", cdsClient.Auth)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := cdsClient.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	log.Printf("BUCKET DELETE RESPONSE STATUS\n%v\n\n", response.Status)
}

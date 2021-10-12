package clientApi

import (
	"contentWorkflow/baseClient"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const host = "http://localhost:30332/client_api/v1/"

type ClientClient struct {
	client baseClient.BaseClient
}

func NewClientClient() ClientClient {
	bc := baseClient.NewBaseClient()
	return ClientClient{
		client: bc,
	}
}

func (c *ClientClient) DownloadEntry(bucketId string, entryId string, versionId string) {
	log.Println("DOWNLOADING ENTRY")
	url := fmt.Sprintf("%sbuckets/%s/entries/%s/versions/%s/content/", host, bucketId, entryId, versionId)
	fmt.Printf("\nDOWNLOAD ENTRY URL  %s\n\n", url)
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Add("Host", "14aa7ae-da40-4b75-8bd2-048c80dd3d93.client-api-e2e.unity3dusercontent.com")
	if err != nil {
		log.Fatalln(err)
	}
	response, err := c.client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}
	defer response.Body.Close()

	fmt.Printf("\nDOWNLOAD ENTRY RESPONSE STATUS %s\n", response.Status)

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("DONWLOAD ENTRY RESPONSE BODY%s", responseBody)
}

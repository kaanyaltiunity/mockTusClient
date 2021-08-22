package baseClient

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type BaseClient struct {
	http.Client
	Auth string
}

func NewBaseClient() BaseClient {
	return BaseClient{
		Client: http.Client{
			Timeout: time.Duration(20 * time.Second),
		},
	}
}

func (baseClient *BaseClient) SetAuth(authType string, authData string) {
	log.Printf("SETTING AUTH\n\"%s\"\n\n", fmt.Sprintf("%s %s", authType, authData))
	baseClient.Auth = fmt.Sprintf("%s %s", authType, authData)
}

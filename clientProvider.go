package main

import (
	"contentWorkflow/cdsClient"
	"contentWorkflow/gatewayClient"
	"fmt"
	"os"
	"strings"
)

const (
	gateway = "gateway"
	cds     = "cds"
)

type ClientProvider struct{}

func (clientProvider ClientProvider) GetClient() (ContentClient, error) {
	switch os.Args[1] {
	case gateway:
		return gatewayClient.NewGatewayClient(getProjectId()), nil
	case cds:
		return cdsClient.NewCdsClient(getProjectId()), nil
	default:
		return nil, fmt.Errorf("Client type is not provided")
	}
}

func getProjectId() string {
	projectIdVar := fmt.Sprintf("PROJECT_ID_%s", strings.ToUpper(os.Getenv("ENV")))
	return os.Getenv(projectIdVar)
}

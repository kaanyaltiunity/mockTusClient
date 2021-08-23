package main

import (
	"contentWorkflow/cdsClient"
	"contentWorkflow/gatewayClient"
	"fmt"
	"os"
)

const (
	gateway = "gateway"
	cds     = "cds"
)

type ClientProvider struct{}

func (clientProvider ClientProvider) GetClient() (ContentClient, error) {
	switch os.Args[1] {
	case gateway:
		return gatewayClient.NewGatewayClient(os.Getenv("PROJECT_ID")), nil
	case cds:
		return cdsClient.NewCdsClient(os.Getenv("PROJECT_ID")), nil
	default:
		return nil, fmt.Errorf("Client type is not provided")
	}
}

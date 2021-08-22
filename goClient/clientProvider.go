package main

import (
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
		fmt.Println("HELLO")
		return nil, nil
	default:
		return nil, fmt.Errorf("Client type is not provided")
	}
}

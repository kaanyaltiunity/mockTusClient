# Dev Overview

This tool is a simple client implementation for some endpoints of the content delivery management api. The client can send requests to both the gateway and non-gateway endpoints. 

In order to send requests to the gateway endpoints run the following command
```
go run . gateway
```

To send requests to the non-gateway endpoints run the following command
```
go run . cds
```

Required environment variables can be found in the `.env.sample` file (also listed below). 
```
BEARER_TOKEN=
API_KEY_STAGING=
API_KEY_DEV=
PROJECT_ID_STAGING=
PROJECT_ID_DEV=
ENV=dev
GO_TUS_ENABLED=false
```

`BEARER_TOKEN` is used for gateway requests. It should be generated using the tooling provided by the service gateway team

`API_KEY_DEV` and `API_KEY_STAGING` are used for non-gateway requests. The staging api key can be found on the unity dashboard. The dev key is the key that's used for e2e tests. 

`GO_TUS_ENABLED` toggles tus client on and off.

When tus uploads are enabled you may want to change the chunk size of each upload. In order to do so please look at the `uploadWithGoTus` function in both the gateway and the cds clients
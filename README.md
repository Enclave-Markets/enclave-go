# enclave-go

<p align="center">
    <a href="https://github.com/Enclave-Markets/enclave-go" alt="enclave go">
        <img src="https://edent.github.io/SuperTinyIcons/images/svg/github.svg" width="50" /></a>
    <a href="https://twitter.com/enclavemarkets" alt="Enclave Twitter">
        <img src="https://edent.github.io/SuperTinyIcons/images/svg/x.svg" width="50"/></a>
    <a href="https://www.enclave.market/" alt="Enclave Market">
        <img src="https://pbs.twimg.com/profile_images/1650572649284931585/rbv_Z4Lr_400x400.jpg" width="50"/></a>
        
</p>

This is and example Go SDK for [Enclave Markets](https://enclave.market/). 

It provides a simple interface for interacting with the spot market [Enclave API](https://docs.enclave.market/).


## Installation

```bash
go get github.com/Enclave-Markets/enclave-go
```

## Usage

```go
package main

import (
	"github.com/Enclave-Markets/enclave-go/api_client"
)

func main() {

	client, err := api_client.NewApiClientFromEnv("sandbox")
	if err != nil {
		return
	}

	client.WaitForEndpoint()
	_, err = client.Hello()
	if err != nil {
		return
	}
}
```

## Examples

An example of interacting with a spot market on Enclave's sandbox environment can be found in `main.go` and can be run using:

```shell
export enclave_key="YOUR_API_KEY"
export enclave_secret="YOUR_API_SECRET"
go run ./...
```

API keys for Enclave's sandbox environment can be found [here](https://sandbox.enclave.market/) by first connecting a wallet and then accessing account settings.

## Support

Supports Go 1.22+
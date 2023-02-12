# AWS Credentials Cache for Go

This module provides credentials caching utilities that are partially compatible with AWS CLI.

## Installation

```shell
go get github.com/Aton-Kish/aws-credscache-go
```

## Usage

```go
package main

import (
	"context"
	"log"

	credscache "github.com/Aton-Kish/aws-credscache-go/sdkv2"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithAssumeRoleCredentialOptions(func(options *stscreds.AssumeRoleOptions) {
		options.TokenProvider = stscreds.StdinTokenProvider
	}))
	if err != nil {
		log.Fatal(err)
	}

	// Inject file cache provider
	if _, err := credscache.InjectFileCacheProvider(&cfg); err != nil {
		log.Fatal(err)
	}

	// client := ec2.NewFromConfig(cfg)
}
```

See [exmples](./_examples/) for more details.

## Development

### Setup

```shell
go mod tidy
```

### Generate code

```shell
rm -rf internal/mock
go generate ./...
```

### Test

```shell
: simple
go test ./...
: verbose
go test -v ./...
```

### Doc

```shell
go run golang.org/x/tools/cmd/godoc@latest -http ":6060"
```

## License

This library is licensed under the MIT License, see [LICENSE](./LICENSE).

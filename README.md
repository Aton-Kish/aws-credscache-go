# AWS Credentials Cache for Go

This module provides credentials caching utilities that are compatible with the AWS CLI.

## Motivation

The AWS SDK has a feature of an in-memory cache for credentials.
However, it doesn't work effectively for use cases of the short-lifespan process like CLI.

![nocache](./_examples/cli/images/gif/sdkv2_nocache.gif)  
An MFA token code will be requested every time.
It's very bothering.

Although the AWS CLI saves credentials into `$HOME/.aws/cli/cache`, the AWS SDK does not support it.
This module provides an easy way to apply a file-caching feature that has compatibility with the AWS CLI.

![cache](./_examples/cli/images/gif/sdkv2_cache.gif)  
![cache shared with AWS CLI](./_examples/cli/images/gif/sdkv2_cache_awscli.gif)  
You will input an MFA token code only once and can also share the cache with the AWS CLI.

See [exmples](./_examples/) for more details.

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

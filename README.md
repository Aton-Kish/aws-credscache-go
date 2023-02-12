# AWS Credentials Cache for Go

## Development

### setup

```shell
go mod tidy
```

### generate code

```shell
rm -rf internal/mock
go generate ./...
```

### test

```shell
: simple
go test ./...
: verbose
go test -v ./...
```

### doc

```shell
go run golang.org/x/tools/cmd/godoc@latest -http ":6060"
```

## License

This library is licensed under the MIT License, see [LICENSE](./LICENSE).

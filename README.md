# go-boilerplate

Go boilerplate project

## Getting started

### Install Dependencies

From the project root, run:

```shell
go build ./...
go test ./...
go mod tidy
```

### Update/Upgrade dependencies

```shell
go get -u ./...
```

### Run dev

```shell
# API server
go run main.go api --port 8000 --shutdown_time 10
```

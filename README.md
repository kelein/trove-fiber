# Project Template

![Go](https://img.shields.io/github/go-mod/go-version/kelein/git?style=for-the-badge&logo=go&logoColor=%2300ADD8&color=%2300ADD8)     [![Helm](https://img.shields.io/badge/helm_v3-%23101683.svg?style=for-the-badge&logo=helm&logoColor=white)](https://github.com/kelein/git)      [![GitHub](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)](https://github.com/kelein/git/actions)      [![Build](https://img.shields.io/github/actions/workflow/status/kelein/git/codeql.yml?style=for-the-badge&logo=github)](https://github.com/kelein/git/actions)

## Install Binary

```erlang
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
```

## Create Service

```erlang
# Create a template project
kratos new server
cd server

# Add a proto template
kratos proto add api/server/server.proto

# Generate the proto code
kratos proto client api/server/server.proto

# Generate the source code of service by proto file
kratos proto server api/server/server.proto -t internal/service

go generate ./...
go build -o ./bin/ ./...
./bin/server -conf ./configs
```

## Makefile Generate

```erlang
❯ make

Usage:
 make [target]

Targets:
init:           init env
config:         generate proto
api:            generate api proto
build:          build binary
generate:       generate
all:            generate all
help:           show help
```

## Automate Initial

```erlang
# install wire
go get github.com/google/wire/cmd/wire

# generate wire
cd cmd/server
wire
```

## Docker Image

```erlang
# build
docker build -t <your-docker-image-name> .

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/configs>:/data/conf <image-name>
```

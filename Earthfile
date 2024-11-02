VERSION 0.8
FROM golang:1.22-bookworm
WORKDIR /workspace

all:
  ARG VERSION=dev
  BUILD --platform=linux/amd64 --platform=linux/arm64 --platform=linux/riscv64 +docker

build:
  ARG GOOS=linux
  ARG GOARCH=amd64
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=1 go build --ldflags '-linkmode external' -o action-mergeusr main.go
  SAVE ARTIFACT ./action-mergeusr AS LOCAL dist/action-mergeusr-${GOOS}-${GOARCH}

tidy:
  LOCALLY
  ENV GOTOOLCHAIN=go1.22.1
  RUN go mod tidy
  RUN go fmt ./...

lint:
  FROM golangci/golangci-lint:v1.61.0
  WORKDIR /workspace
  COPY . ./
  RUN golangci-lint run --timeout 5m ./...

test:
  FROM +tools
  ARG TARGETARCH
  COPY +build/immutos ./dist/immutos-linux-${TARGETARCH}
  COPY . ./
  WITH DOCKER
    RUN go test -coverprofile=coverage.out -v ./...
  END
  SAVE ARTIFACT ./coverage.out AS LOCAL coverage.out

docker:
  FROM debian:trixie-slim
  RUN apt update \
    && apt install -y fakechroot \
    && rm -rf /var/lib/apt/lists/*
  COPY LICENSE /usr/share/doc/action-mergeusr/copyright
  ARG TARGETARCH
  COPY (+build/action-mergeusr --GOARCH=$TARGETARCH) /usr/local/bin/action-mergeusr
  ENTRYPOINT ["/usr/local/bin/action-mergeusr"]
  ARG VERSION=dev
  SAVE IMAGE --push ghcr.io/immutos/action-mergeusr:${VERSION}
  SAVE IMAGE --push ghcr.io/immutos/action-mergeusr:latest
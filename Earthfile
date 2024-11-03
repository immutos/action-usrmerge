VERSION 0.8

all:
  ARG VERSION=dev
  BUILD --platform=linux/amd64 --platform=linux/arm64 +docker

build:
  FROM golang:1.22-bookworm
  WORKDIR /workspace
  RUN apt update
  RUN apt install -y gcc-aarch64-linux-gnu gcc-x86-64-linux-gnu
  ARG GOOS=linux
  ARG GOARCH=amd64
  ENV TARGET_TRIPLET=$(echo "$GOARCH" | sed -e 's/amd64/x86_64-linux-gnu/' -e 's/arm64/aarch64-linux-gnu/')
  ENV CC=$(echo "$TARGET_TRIPLET-gcc")
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
  FROM golang:1.22-bookworm
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN go test -coverprofile=coverage.out -v ./...
  SAVE ARTIFACT ./coverage.out AS LOCAL coverage.out

docker:
  FROM debian:bookworm-slim
  RUN apt update \
    && apt install -y fakechroot \
    && rm -rf /var/lib/apt/lists/*
  COPY LICENSE /usr/share/doc/action-mergeusr/copyright
  ARG NATIVEARCH
  ARG TARGETARCH
  COPY --platform=linux/$NATIVEARCH (+build/action-mergeusr --GOARCH=$TARGETARCH) /usr/local/bin/action-mergeusr
  ENTRYPOINT ["/usr/local/bin/action-mergeusr"]
  ARG VERSION=dev
  SAVE IMAGE --push ghcr.io/immutos/action-mergeusr:${VERSION}-bookworm
  SAVE IMAGE --push ghcr.io/immutos/action-mergeusr:bookworm
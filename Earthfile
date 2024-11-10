VERSION 0.8

all:
  ARG VERSION=dev
  BUILD --platform=linux/amd64 --platform=linux/arm64 +docker-bookworm
  BUILD --platform=linux/amd64 --platform=linux/arm64 --platform=linux/riscv64 +docker-trixie

build-bookworm:
  FROM debian:bookworm-slim
  WORKDIR /workspace
  RUN apt update
  RUN apt install -y ca-certificates golang-go \
    gcc-x86-64-linux-gnu gcc-aarch64-linux-gnu
  ARG GOOS=linux
  ARG GOARCH=amd64
  ENV TARGET_TRIPLET=$(echo "$GOARCH" | sed -e 's/amd64/x86_64-linux-gnu/' -e 's/arm64/aarch64-linux-gnu/')
  ENV CC=$(echo "$TARGET_TRIPLET-gcc")
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=1 go build --ldflags '-linkmode external' -o action-usrmerge main.go
  SAVE ARTIFACT ./action-usrmerge AS LOCAL dist/action-usrmerge-${GOOS}-${GOARCH}

build-trixie:
  FROM debian:trixie-slim
  WORKDIR /workspace
  RUN apt update
  RUN apt install -y ca-certificates golang-go \
    gcc-x86-64-linux-gnu gcc-aarch64-linux-gnu gcc-riscv64-linux-gnu
  ARG GOOS=linux
  ARG GOARCH=amd64
  ENV TARGET_TRIPLET=$(echo "$GOARCH" | sed -e 's/amd64/x86_64-linux-gnu/' -e 's/arm64/aarch64-linux-gnu/' -e 's/riscv64/riscv64-linux-gnu/')
  ENV CC=$(echo "$TARGET_TRIPLET-gcc")
  COPY go.mod go.sum ./
  RUN go mod download
  COPY . .
  RUN CGO_ENABLED=1 go build --ldflags '-linkmode external' -o action-usrmerge main.go
  SAVE ARTIFACT ./action-usrmerge AS LOCAL dist/action-usrmerge-${GOOS}-${GOARCH}

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

docker-bookworm:
  FROM debian:bookworm-slim
  RUN apt update \
    && apt install -y fakechroot \
    && rm -rf /var/lib/apt/lists/*
  COPY LICENSE /usr/share/doc/action-usrmerge/copyright
  ARG NATIVEARCH
  ARG TARGETARCH
  COPY --platform=linux/$NATIVEARCH (+build-bookworm/action-usrmerge --GOARCH=$TARGETARCH) /usr/local/bin/action-usrmerge
  ENTRYPOINT ["/usr/local/bin/action-usrmerge"]
  ARG VERSION=dev
  SAVE IMAGE --push ghcr.io/immutos/action-usrmerge:${VERSION}-bookworm
  SAVE IMAGE --push ghcr.io/immutos/action-usrmerge:bookworm

docker-trixie:
  FROM debian:trixie-slim
  RUN apt update \
    && apt install -y fakechroot \
    && rm -rf /var/lib/apt/lists/*
  COPY LICENSE /usr/share/doc/action-usrmerge/copyright
  ARG NATIVEARCH
  ARG TARGETARCH
  COPY --platform=linux/$NATIVEARCH (+build-trixie/action-usrmerge --GOARCH=$TARGETARCH) /usr/local/bin/action-usrmerge
  ENTRYPOINT ["/usr/local/bin/action-usrmerge"]
  ARG VERSION=dev
  SAVE IMAGE --push ghcr.io/immutos/action-usrmerge:${VERSION}-trixie
  SAVE IMAGE --push ghcr.io/immutos/action-usrmerge:trixie
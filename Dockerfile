FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.22.4 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /src

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-w -s" -o ./bin/gogl cmd/cli/main.go

FROM --platform=${TARGETPLATFORM:-linux/amd64} debian:bookworm-20240612-slim

LABEL org.opencontainers.image.source https://github.com/WatcherWhale/gogl-ci

WORKDIR /

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates

RUN mkdir -p /root/.config/gogl

COPY --from=build /src/bin/gogl /usr/bin/gogl

ENTRYPOINT ["/usr/bin/gogl"]

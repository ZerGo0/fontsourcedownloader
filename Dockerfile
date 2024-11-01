# syntax = docker/dockerfile:1.4

# get modules, if they don't change the cache can be used for faster builds
FROM golang:1.23.1@sha256:4f063a24d429510e512cc730c3330292ff49f3ade3ae79bda8f84a24fa25ecb0 AS base
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /src
COPY go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# build th application
FROM base AS build
# temp mount all files instead of loading into image with COPY
# temp mount module cache
# temp mount go build cache
RUN --mount=target=. \
    --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -ldflags="-w -s" -o /app/main ./main.go

# Import the binary from build stage

FROM gcr.io/distroless/static:nonroot@sha256:3a03fc0826340c7deb82d4755ca391bef5adcedb8892e58412e1a6008199fa91 as prd
COPY --link --from=build /app/main /
# this is the numeric version of user nonroot:nonroot to check runAsNonRoot in kubernetes
USER 65532:65532
ENTRYPOINT ["/main"]

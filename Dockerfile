# syntax=docker/dockerfile:1.3
ARG GO_VERSION=1.17

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM pratikimprowise/upx:3.96 AS upx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
ENV GO111MODULE=auto
RUN apk --update add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

## bin
FROM vendored AS bin
ENV CGO_ENABLED=0
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="zipper" \
    --main="./cmd/zipper" \
    --dist="/out" \
    --artifacts="bin" \
    --artifacts="archive" \
    --snapshot="no"

## Slim bin
FROM vendored AS bin-slim
ENV CGO_ENABLED=0
COPY --from=upx / /
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name="zipper-slim" \
    --flags="-trimpath" \
    --ldflags="-s -w" \
    --main="./cmd/zipper" \
    --dist="/out" \
    --artifacts="bin" \
    --artifacts="archive" \
    --snapshot="no" \
    --post-hooks="upx -v --ultra-brute --best -o /usr/local/bin/{{ .ProjectName }}{{ .Ext }}"

## get binary out
### non slim binary
FROM scratch AS artifact
COPY --from=bin /out /
###

### slim binary
FROM scratch AS artifact-slim
COPY --from=bin-slim /out /
###

### All binaries
FROM scratch AS artifact-all
COPY --from=bin /out /
COPY --from=bin-slim /out /
###
##

### Testing
FROM vendored as test
COPY . .
ARG TARGETPLATFORM
RUN apk --update add --no-cache gcc
RUN CGO_ENABLED=0 go test -v ./...
RUN go test -v -race ./...

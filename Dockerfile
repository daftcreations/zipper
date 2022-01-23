# syntax=docker/dockerfile:1.3

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM pratikimprowise/upx:3.96 AS upx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
ENV CGO_ENABLED=0
ENV GO111MODULE=auto
RUN apk --update add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,target=.,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored as test
RUN go test -v ./..
RUN go test -v race ./...

## bin
FROM vendored AS bin
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
    --post-hooks="sh -c 'upx -v --ultra-brute --best -o /usr/local/bin/{{ .ProjectName }}{{ .Ext }} || true'"

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

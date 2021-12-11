ARG GO_VERSION=1.17

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "1001" \
    "appuser"
RUN apk add --no-cache ca-certificates gcc file git linux-headers musl-dev tar
WORKDIR /src

FROM base AS build
ARG TARGETPLATFORM
ARG GIT_REF
RUN --mount=type=bind,target=/src,rw \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=target=/go/pkg/mod,type=cache \
  goreleaser-xx --debug \
    --name "zipper" \
    --dist "/out" \
    --hooks="go mod tidy" \
    --hooks="go mod download" \
    --main="." \
    --ldflags="-s -w" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifacts
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

# FROM --platform=$BUILDPLATFORM pratikimprowise/golang:1.17 as build2
# COPY --from=build /usr/local/bin/zipper /usr/local/bin/zipper
# RUN  strip /usr/local/bin/zipper && \
#   /usr/local/bin/upx -9 /usr/local/bin/zipper

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /usr/local/bin/zipper /usr/local/bin/zipper
USER appuser:appuser
WORKDIR /tmp
WORKDIR /
ENTRYPOINT [ "/usr/local/bin/zipper" ]

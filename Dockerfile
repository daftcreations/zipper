ARG GO_VERSION=1.17

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:1.1.0 AS goreleaser-xx
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
RUN apk add --no-cache git
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
    --build-pre-hooks="go mod tidy" \
    --build-pre-hooks="go mod download" \
    --main="." \
    --ldflags="-s -w" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifacts
COPY --from=build /out/*.tar.gz /
COPY --from=build /out/*.zip /

FROM --platform=$BUILDPLATFORM pratikimprowise/golang:${GO_VERSION} as build2
ARG TARGETPLATFORM
COPY --from=build /usr/local/bin/zipper /usr/local/bin/zipper
RUN upx -9 /usr/local/bin/zipper

FROM scratch
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=build2 /usr/local/bin/zipper /usr/local/bin/zipper
USER appuser:appuser
WORKDIR /tmp
WORKDIR /
ENTRYPOINT [ "/usr/local/bin/zipper" ]

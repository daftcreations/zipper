ARG GO_VERSION=1.17

FROM --platform=$BUILDPLATFORM crazymax/goreleaser-xx:latest AS goreleaser-xx
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine AS base
COPY --from=goreleaser-xx / /
ENV CGO_ENABLED=0
RUN apk add --no-cache git
WORKDIR /src

FROM base AS vendored
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/go/pkg/mod \
  go mod tidy && go mod download

FROM vendored AS zipper
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name "zipper" \
    --ldflags="-s -w" \
    --flags="-trimpath" \
    --dist "/out" \
    --main="cmd/zipper/" \
    --files="LICENSE" \
    --files="README.md"

FROM vendored AS resizer
ARG TARGETPLATFORM
RUN --mount=type=bind,source=.,target=/src,rw \
  --mount=type=cache,target=/root/.cache \
  --mount=type=cache,target=/go/pkg/mod \
  goreleaser-xx --debug \
    --name "resizer" \
    --ldflags="-s -w" \
    --flags="-trimpath" \
    --dist "/out" \
    --main="cmd/resizer/" \
    --files="LICENSE" \
    --files="README.md"

FROM scratch AS artifacts
COPY --from=zipper /out/*.tar.gz /out/*.zip  /
COPY --from=resizer /out/*.tar.gz /out/*.zip  /

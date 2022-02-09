# Zipper [![Go Report Card](https://goreportcard.com/badge/github.com/pratikbalar/zipper)](https://goreportcard.com/report/github.com/pratikbalar/zipper)

<img src="docs/zipper.png" alt="zipper logo" width="100" height="100"/>

`zipper` create multiple zip files of X MB

> for my friend savan

![Views](https://dynamic-badges.maxalpha.repl.co/views?id=pratikbalar.zipper&style=for-the-badge&color=black)

## Usage

```bash
zipper <size> <target>
```

- Size in KB. Default 3000 which is 3MB
- target: folder/directory to compress

## Local

> Create docker buildx builder if using first time
> ```docker buildx create --use```

> `slim` variant is UPX compressed

```shell
git clone --depth 1 https://github.com/pratikbalar/zipper.git zipper
cd zipper

## Get zipper binary
docker buildx bake

## Get slim zipper binary
docker buildx bake artifact-slim

## Get slim and standard multi-platform zipper binaries
docker buildx bake artifact-all
```

*Check `dist` directory for binaries*

```shell
## Build image
docker buildx bake image

## Build slim image
docker buildx bake image-slim

## Build multi-platform image
docker buildx bake image-all

## Build multi-platform slim image
docker buildx bake image-all-slim
```

Container image availabel for

- **linux**: `amd64`, `386`, `arm64`, `riscv64`, `ppc64le`, `s390x`, `mips64le`, `mips64`, `arm/v7`, `arm/v6`

- **windows**: *available soon*

Image available for

- **linux**: `amd64`, `386`, `arm64`, `riscv64`, `ppc64le`, `s390x`, `mips64le`, `mips64`, `arm/v7`, `arm/v6`

- **darwin**: `amd64`(Intel), `arm64` (M1)

- **windows**: `amd64`, `386`, `arm64`, `arm`

- **freebsd**: `amd64`, `386`, `arm64`, `arm`

## Debug

### Testing

```shell
go test -v ./...
go test -v -race ./...
```

### Profiling

```shell
go build -ldflags="-X main.profEnable=true" ./cmd/zipper/
./zipper 10000 /home/pratik/workspace/pratikbalar/zipper/test
```

#### CPU and Memory profiling

```shell
go tool pprof -http=:8080 mem.pprof
```

<!-- **OR** -->
<!--
```shell
go test -cpuprofile cpu.prof -memprofile mem.prof -bench ./cmd/zipper/
``` -->

#### Tracing

```shell
go tool trace trace.out
```

> <div>Icons made by <a href="https://www.freepik.com" **title**="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>

## To-Do

- [ ] Use existing buffer for zip
- [ ] Replace queue with channels
- [ ] Create zipping worker with channels
- [ ] Add structured logging `<TYPE>: <TIME> :`
- [x] Check goroutine leak
- [x] Don't change source dir name
- [x] Configure size of zip from ~~env~~ cli
- [x] docs: mkdocs CI
- [x] What if photo size is less then given zip size
- [x] Go routines to handle zipping
- [x] Testing
  - [x] e2e
- [x] **README** Badges

---

*May the source be with you*

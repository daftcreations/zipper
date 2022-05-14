# Zipper

<div align="center">
    <img src="docs/zipper.png" style="width:200px; border-radius:10px;"/>
</div><br>
<p align="center">
  <a href="https://github.com/daftcreations/zipper/blob/master/LICENSE"><img src="https://img.shields.io/github/license/daftcreations/zipper?style=flat-square"/></a>
  <a href="https://github.com/daftcreations/zipper/actions/workflows/release.yml"><img src="https://img.shields.io/github/workflow/status/daftcreation/zipper/goreleaser?style=flat-square"/></a>
  <a href="https://github.com/daftcreations/zipper/releases"><img src="https://img.shields.io/github/v/tag/daftcreations/zipper?style=flat-square"/></a>
  <a href=""><img src="https://goreportcard.com/badge/github.com/pratikbin/zipper?label=Lines%20of%20code&style=flat-square"/></a>
  <a href="https://goreportcard.com/report/github.com/pratikbin/zipper"><img src="https://img.shields.io/tokei/lines/github/daftcreations/zipper?label=Lines%20of%20code&style=flat-square"/></a>
  <a href="https://discord.com/channels/960581263264219186/960618259244257330"><img src="https://img.shields.io/discord/960581263264219186?label=%20&logo=discord&style=flat-square"/></a>
</p>

`zipper` creat2e multiple zip files of X MB

## Usage

```shell
zipper <size> # size in KiloBytes
Splitting into <size> KB

Enter path you want to zip:
>
```

## Local

> Create docker buildx builder if using first time
> ```docker buildx create --use```
>
> `slim` variant is UPX compressed

```shell
git clone --depth 1 https://github.com/pratikbin/zipper.git zipper
cd zipper

## Get zipper binary
docker buildx bake

## Get slim zipper binary
docker buildx bake artifact-slim

## Get slim and standard multi-platform zipper binaries
docker buildx bake artifact-all
```

> Check `dist` directory for binaries*

Binaries available for

- linux: `amd64`, `386`, `arm64`, `riscv64`, `ppc64le`, `s390x`, `mips64le`, `mips64`, `arm/v7`, `arm/v6`

- darwin: `amd64`(Intel), `arm64` (M1)

- windows: `amd64`, `386`, `arm64`, `arm`

- freebsd: `amd64`, `386`, `arm64`, `arm`

## Debug

## Testing

```shell
go test -v -coverprofile cover.out .
go tool cover -html=cover.out -o cover.html
open cover.html
```

## Profiling

```shell
go build -ldflags="-X main.profEnable=true" ./cmd/zipper/
```

### CPU and Memory profiling

```shell
go tool pprof -http=:8080 mem.pprof &
go tool pprof -http=:8081 cpu.pprof
```

<!-- **OR** -->
<!--
```shell
go test -cpuprofile cpu.prof -memprofile mem.prof -bench ./cmd/zipper/
``` -->

### Tracing

```shell
go tool trace trace.out
```

> <div>Icons made by <a href="https://www.freepik.com" **title**="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>

## TODO

- [ ] Use existing buffer for zip
- [ ] Replace queue with channels
- [ ] Add structured logging `<TYPE>: <TIME> :`
- [x] Check goroutine leak
- [x] Create zipping worker with channels
- [x] Go routines to handle zipping
- [x] Don't change source dir name
- [x] Configure size of zip from ~~env~~ cli
- [x] docs: mkdocs CI
<!-- Was part of resizer CLI - [x] What if photo size is less then given zip size -->
- [x] Testing
  - [x] e2e
- [x] **README** Badges

---

<div align="center">
    <a href="https://discord.com/channels/960581263264219186/960618259244257330"><img src="https://img.shields.io/badge/Discord-5865F2?style=for-the-badge&logo=discord&logoColor=white"/></a>
    <a href="https://www.youtube.com/c/DaftCreation/playlists"><img src="https://img.shields.io/youtube/channel/subscribers/UCDrfHGsm6bJI7Sli7vlcteA?label=YT&logo=youtube&style=for-the-badge"/></a>
    <a href="https://twitter.com/daftcreations"><img src="https://img.shields.io/twitter/follow/daftcreation?logo=twitter&style=for-the-badge"/></a>
    <a href="https://www.instagram.com/daft.creations/"><img src="https://img.shields.io/badge/Instagram-E4405F?style=for-the-badge&logo=instagram&logoColor=white"/></a>
</div>

<div align="center">
    <div style="display:flex; justify-content:space-around;">
        <h3 style="margin:-5px 10px 5px;">Contributors</h3>
        <hr align="left" width="20%">
    </div>
    <img src="https://contrib.rocks/image?repo=daftcreations/zipper&columns=80" style="width:150px;"/>
</div>

### Stargazers over time

<center>
    <a href="https://starchart.cc/daftcreations/zipper"><img src="https://starchart.cc/daftcreations/zipper.svg" width="80%"/></a>
</center>

> *May the source be with you*

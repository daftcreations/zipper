# Zipper

<img src="docs/zipper.png" alt="zipper logo" width="100" height="100"/>

`zipper` create multiple zip files of less then X MB

> for my friend savan

![Views](https://dynamic-badges.maxalpha.repl.co/views?id=pratikbalar.zipper&style=for-the-badge&color=black)

## Usage

```bash
zipper <(Optional) size> <(Optional) target>
```

- (Optional) Size in KB. Default 5000 which is 5MB, e.g. `zipper 5000`
- (Optional) target: folder/directory which going to compress it will ask you in
prompt it you won't add target in args

### Binary

1. Download `zipper` as per your OS from [releases](https://github.com/pratikbalar/zipper/releases) page
page and extract it

    Binaries available for fllowing OS and architecture

    - darwin
        - amd64 (Mac Intel)
        - arm64 (Mac M1 All)
    - windows
        - 386
        - amd64
        - arm6
    - linux
        - 386
        - amd64
        - arm-v5
        - arm-v6
        - arm-v7
        - arm64
        - ppc64le
        - riscv64
        - s390x

2. Open terminal (cmd in case of windows) in extracted path

    linux & Mac

    ```
    ./zipper 5000 \home\user\test
    ```

    Windows

    ```
    .\zipper.exe 5000 D:\test
    ```

### Container

Container image available for fllowing architecture

- linux 386
	- amd64
	- arm-v6
	- arm-v7
	- arm64
	- ppc64le


1. Go to folder/directory you want to zip

2. Now run

```bash
docker run -v $(pwd):/zip -it pratikimprowised/zipper 5000 /zip/test
```

> <div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>

## To-Do

- [x] Drag and drop based for windows
- [x] Don't change file name
- [x] Configure size of zip from env
- [x] docs: mkdocs CI
- [ ] UPX binaries for multiplateform containers
- [ ] What if dir size is less then given zip size
- [ ] Testing
  - [ ] unit
  - [ ] integration wrt multi os
  - [ ] e2e
- [ ] `UPX`ed binaries in releases
- [ ] README Badges

# Zipper

<img src="docs/zipper.png" alt="zipper logo" width="100" height="100"/>

`zipper` create multiple zip files of less then X MB

> for my friend savan

![Views](https://dynamic-badges.maxalpha.repl.co/views?id=pratikbalar.zipper&style=for-the-badge&color=black)

## Usage

```bash
zipper <(Optional) size> <(Optional) target>
```

- (Optional) Size in KB. Default 3000 which is 3MB, e.g. `zipper 3000`
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
    ./zipper 3000 \home\user\test
    ```

    Windows

    ```
    .\zipper.exe 3000 D:\test
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
docker run -v $(pwd):/zip -it pratikimprowised/zipper 3000 /zip/test
```

## Setup zipper

### Linux & Mac

- Open terminal in extracted zipper folder/directory

- Move `zipper` to bin

```bash
mv zipper /usr/local/bin/zipper
```

- Now `zipper` should be available user wide

### Windows

- Open powershell in extracted zipper folder/directory

- Move `zipper` to `c:\bins\`

```
New-Item -Path "c:\" -Name "bins" -ItemType "directory"
Move-Item -Path zipper -Destination c:\bins\
```

- Add `c:\bins\` to user's `PATH`

```
$PATH = [Environment]::GetEnvironmentVariable("PATH")
$bins_path = "c:\bins\"
[Environment]::SetEnvironmentVariable("PATH", "$PATH;$bins_path")
```

- Now `zipper` should be available user wide

## To-Do

- [x] Drag and drop based for windows
- [x] Don't change file name
- [x] Configure size of zip from env
- [x] docs: mkdocs CI
- [x] UPX binaries for multiplateform containers
- [ ] `UPX`ed binaries in releases
- [ ] What if dir size is less then given zip size
- [ ] Go routines to handle zipping
- [ ] Testing
  - [x] unit
  - [ ] integration wrt multi os
  - [ ] e2e
- [ ] README Badges

> <div>Icons made by <a href="https://www.freepik.com" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>

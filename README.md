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

## To-Do

- [x] Don't change source dir name
- [x] Add dirname in created zip
- [x] Configure size of zip from ~~env~~ cli
- [x] docs: mkdocs CI
- [x] What if photo size is less then given zip size
- [x] Go routines to handle zipping
- [x] Testing
  - [x] unit
  - [x] integration wrt multi os
  - [x] e2e
- [ ] **README** Badges

> <div>Icons made by <a href="https://www.freepik.com" **title**="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon">www.flaticon.com</a></div>

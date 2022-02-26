# Bazel

## Update

### Build files

```shell
bazel run //:gazelle
```

### Go deps

```shell
bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=deps.bzl%go_dependencies
```

variable "GO_VERSION" {
  default = "1.17"
}

variable "REPO" {
  default = "pratikbalar/zipper"
}

variable "VERSION" {
  default = "edge"
}

variable "GIT_SHA" {
  default = "0000000"
}

variable "TAGS" {
  default = [
    "${REPO}:latest",
    "${REPO}:${VERSION}",
    "${REPO}:${VERSION}-${GIT_SHA}"
  ]
}

variable "TAGS_SLIM" {
  default = [
    "${REPO}:slim",
    "${REPO}:${VERSION}-slim",
    "${REPO}:${VERSION}-${GIT_SHA}-slim"
  ]
}


target "_common" {
  args = {
    GO_VERSION = GO_VERSION
  }
}

target "_labels" {
  labels = {
    "org.opencontainers.image.title"         = "zipper",
    "org.opencontainers.image.base.name "    = "scratch",
    "org.opencontainers.image.licenses"      = "Apache-2.0",
    "org.opencontainers.image.description"   = "Multiple zip files of less then X MB",
    "org.opencontainers.image.version"       = "${VERSION}",
    "org.opencontainers.image.revision"      = "${GIT_SHA}",
    "org.opencontainers.image.source"        = "https://github.com/pratikbalar/zipper",
    "org.opencontainers.image.documentation" = "https://github.com/pratikbalar/zipper",
  }
}

target "_slim" {
  target = "slim"
  tags   = TAGS_SLIM
}

target "_standard" {
  target = "standard"
  tags   = TAGS
}

target "artifacs" {
  output = ["./dist"]
}

target "image-platform" {
  platforms = [
    "linux/amd64",
    "linux/arm64",
    "linux/riscv64",
    "linux/ppc64le",
    "linux/s390x",
    "linux/386",
    "linux/mips64le",
    "linux/mips64",
    "linux/arm/v7",
    "linux/arm/v6"
  ]
}

target "bin-platform" {
  platforms = [
    "linux/amd64",
    "linux/386",
    "linux/arm64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x",
    "freebsd/amd64",
    "freebsd/386",
    "freebsd/arm64",
    "freebsd/arm",
    "windows/amd64",
    "windows/arm64",
    "windows/arm",
    "windows/386",
    "darwin/amd64",
    "darwin/arm64",
  ]
}

group "default" {
  targets = ["image"]
}

# test
target "test" {
  inherits = ["_common", "image-platform"]
  target   = "test"
}

# Creating standard container image for local docker
target "image" {
  inherits = ["_common", "_standard", "_labels"]
  output   = ["type=docker"]
}

# Creating slim container image for local docker
target "image-slim" {
  inherits = ["_common", "_slim", "_labels"]
  output   = ["type=docker"]
}

# Creating standard container image for all platforms
target "image-all" {
  inherits = ["_common", "image-platform", "_standard", "_labels"]
  target   = "standard"
}

# Creating slim container image for all platforms
target "image-all-slim" {
  inherits = ["_common", "image-platform", "_slim", "_labels"]
  target   = "slim"
}

# Creating all standard artifact for all platforms
target "artifact" {
  inherits = ["_common", "artifacs"]
  target   = "artifact"
}

# Creating all slim artifact for all platforms
target "artifact-slim" {
  inherits = ["_common", "artifacs"]
  target   = "artifact-slim"
}

# Creating all full, slim artifact for all platforms
target "artifact-all" {
  inherits = ["artifact-all", "artifacs", "bin-platform"]
  target   = "artifact-all"
}

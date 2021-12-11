variable "GO_VERSION" {
  default = "1.17"
}
variable "PROTOC_VERSION" {
  default = "3.17.3"
}
variable "GITHUB_REF" {
  default = ""
}
variable "DOCKERHUB_SLUG" {
  default = "pratikimprowised/zipper"
}

// Special target: https://github.com/docker/metadata-action#bake-definition
target "docker-metadata-action" {
  tags = ["${DOCKERHUB_SLUG}:local"]
}

target "_common" {
  args = {
    GO_VERSION = GO_VERSION
    GIT_REF = GITHUB_REF
    BUILDKIT_CONTEXT_KEEP_GIT_DIR = 1
  }
}

group "default" {
  targets = ["image-local"]
}

target "artifact" {
  inherits = ["_common"]
  target = "artifacts"
  output = ["./dist"]
}

target "artifact-all" {
  inherits = ["artifact"]
  platforms = [
    "darwin/amd64",
    "darwin/arm64",
    "linux/386",
    "linux/amd64",
    "linux/arm/v5",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
    "linux/ppc64le",
    "linux/riscv64",
    "linux/s390x",
    "windows/386",
    "windows/amd64",
    "windows/arm64"
  ]
}

target "image" {
  inherits = ["_common", "docker-metadata-action"]
}

target "image-local" {
  inherits = ["image"]
  output = ["type=docker"]
}

target "image-all" {
  inherits = ["image"]
  platforms = [
    "linux/386",
    "linux/amd64",
    "linux/arm/v6",
    "linux/arm/v7",
    "linux/arm64",
  ]
}

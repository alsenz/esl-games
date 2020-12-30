load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_go//go:def.bzl", "go_binary")

# gazelle:prefix github.com/example/project
gazelle(name = "gazelle")

go_binary(
    name = "gotesting",
    srcs = ["src/main.go"],
)

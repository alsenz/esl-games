load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/alsenz/esl-games
gazelle(name = "gazelle")

sh_binary(
    name = "go-get",
    srcs = ["scripts/go-get.sh"],
    data = [
        "deps.bzl",
        "go.mod",
        "//:gazelle-runner",
    ],
    tags = [
        "local",
        "no-sandbox",
    ],
)

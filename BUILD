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

exports_files([
    "WORKSPACE",
    "deps.bzl",
])

#genrule(
#    name = "gen-deepcopy-objects",
#    srcs = [
#        "WORKSPACE",
#        "deps.bzl",
#        "//api/v1alpha1:lesson_types.go",
#    ],
#    outs = ["zz_generated.deepcopy.go"],
#    #cmd = "$(execpath @io_k8s_sigs_controller_tools//cmd/controller-gen:controller-gen) object paths=$< output:dir=$(@D)",
#    cmd = "bazel run @io_k8s_sigs_controller_tools//cmd/controller-gen:controller-gen object paths=$< output:dir=$(@D)",
#    #tags = [
#    #    "local",
#    #"no-sandbox",
#    #],
#    tools = ["@io_k8s_sigs_controller_tools//cmd/controller-gen"],
#)

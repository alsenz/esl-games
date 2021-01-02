load("@com_adobe_rules_gitops//gitops:defs.bzl", "k8s_deploy")
load("@io_bazel_rules_go//go:def.bzl", "go_path")

CONFIGURATION_TABLE = {
    # Configuration Name: Cluster, Primary Namespace, User, Registry
    "microk8s": ("microk8s-cluster", "esl-games", "admin", "localhost:32000"),
}

DEP_MODES = {
    "deps": True,
    "no-deps": False,
}

genrule(
    name = "generate-postgres-secrets",
    srcs = ["manifests/postgres-secret.yaml.tpl"],
    outs = ["manifests/postgres-secret.yaml"],
    cmd = "$(execpath //k8s/scripts:generate-secrets) --source $< --dest $@",
    tools = ["//k8s/scripts:generate-secrets"],
)

go_path(
    name = "controller-gen-go-path",
    visibility = ["//visibility:public"],
    deps = [
        "//api/v1alpha1",
    ],
)

# Note: this genrule could do with being parameterised on the import path - /esl-games/api/v1alpha1 feels a little hard coded to me.
genrule(
    name = "controller-gen-crd",
    srcs = [
        "//api/v1alpha1:v1apha1-srcs",
    ],
    outs = ["manifests/esl-games.alsenz.github.com_lessons.yaml"],
    cmd = "GOCACHE=$$(realpath $(@D))/.gocache;" +
          " GOPATH=$$(realpath $(execpath :controller-gen-go-path));" +
          " OUTDIR=$$(realpath $(@D));" +
          " CMDPATH=$$(realpath $(execpath @io_k8s_sigs_controller_tools//cmd/controller-gen:controller-gen));" +
          " cd $$GOPATH;" +
          " >&2 pwd;" +
          " GOCACHE=$$GOCACHE GOPATH=$$GOPATH $$CMDPATH crd paths=./src/github.com/alsenz/esl-games/api/v1alpha1/ output:crd:dir=$$OUTDIR",
    tags = [
        "no-sandbox",
    ],
    tools = [
        ":controller-gen-go-path",
        "@io_k8s_sigs_controller_tools//cmd/controller-gen",
    ],
)

[
    k8s_deploy(
        name = "lesson-crd." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        gitops = False,
        image_registry = REGISTRY,
        manifests = ["manifests/esl-games.alsenz.github.com_lessons.yaml"],
        namespace = PRIMARY_NAMESPACE,
        user = USER,
    )
    for DEP_MODE in DEP_MODES
    for CONFIGURATION_NAME, (CLUSTER, PRIMARY_NAMESPACE, USER, REGISTRY) in CONFIGURATION_TABLE.items()
]

[
    k8s_deploy(
        name = "namespace." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        gitops = False,
        image_registry = REGISTRY,
        manifests = ["manifests/namespace.yaml"],
        namespace = PRIMARY_NAMESPACE,
        user = USER,
    )
    for DEP_MODE in DEP_MODES
    for CONFIGURATION_NAME, (CLUSTER, PRIMARY_NAMESPACE, USER, REGISTRY) in CONFIGURATION_TABLE.items()
]

[
    k8s_deploy(
        name = "postgres-secret." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        gitops = False,
        image_registry = REGISTRY,
        manifests = [
            "manifests/postgres-secret.yaml",
        ],
        namespace = PRIMARY_NAMESPACE,
        objects = [":namespace.deps." + CONFIGURATION_NAME] if DEP_MODES[DEP_MODE] else [],
        user = USER,
    )
    for DEP_MODE in DEP_MODES
    for CONFIGURATION_NAME, (CLUSTER, PRIMARY_NAMESPACE, USER, REGISTRY) in CONFIGURATION_TABLE.items()
]

[
    k8s_deploy(
        name = "postgres." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        gitops = False,
        image_registry = REGISTRY,
        manifests = [
            "manifests/postgres.yaml",
        ],
        namespace = PRIMARY_NAMESPACE,
        objects = [
            ":namespace.deps." + CONFIGURATION_NAME,
            ":postgres-secret.deps." + CONFIGURATION_NAME,
        ] if DEP_MODES[DEP_MODE] else [],
        user = USER,
    )
    for DEP_MODE in DEP_MODES
    for CONFIGURATION_NAME, (CLUSTER, PRIMARY_NAMESPACE, USER, REGISTRY) in CONFIGURATION_TABLE.items()
]

[
    k8s_deploy(
        name = "matcher." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        gitops = False,
        image_digest_tag = True,
        image_registry = REGISTRY,
        image_repository = "esl-games-matcher",
        images = {
            "matcher-image": "@webapp//:webapp-image",
        },
        manifests = ["manifests/matcher.yaml"],
        namespace = PRIMARY_NAMESPACE,
        objects = [
            ":namespace.deps." + CONFIGURATION_NAME,
            ":postgres-secret.deps." + CONFIGURATION_NAME,
        ] if DEP_MODES[DEP_MODE] else [],
        user = USER,
    )
    for DEP_MODE in DEP_MODES
    for CONFIGURATION_NAME, (CLUSTER, PRIMARY_NAMESPACE, USER, REGISTRY) in CONFIGURATION_TABLE.items()
]

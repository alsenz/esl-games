load("@com_adobe_rules_gitops//gitops:defs.bzl", "k8s_deploy")
load("@com_adobe_rules_gitops//skylib:push.bzl", "k8s_container_push")
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
        "//pkg/k8s/api/v1alpha1",
    ],
)

# Note: this genrule could do with being parameterised on the import path - /esl-games/api/v1alpha1 feels a little hard coded to me.
genrule(
    name = "controller-gen-crd",
    srcs = [
        "//pkg/k8s/api/v1alpha1:v1apha1-srcs",
    ],
    outs = ["manifests/esl-games.alsenz.github.com_lessons.yaml"],
    cmd = "GOCACHE=$$(realpath $(@D))/.gocache;" +
          " GOPATH=$$(realpath $(execpath :controller-gen-go-path));" +
          " OUTDIR=$$(realpath $(@D));" +
          " CMDPATH=$$(realpath $(execpath @io_k8s_sigs_controller_tools//cmd/controller-gen:controller-gen));" +
          " cd $$GOPATH;" +
          " >&2 pwd;" +
          " GOCACHE=$$GOCACHE GOPATH=$$GOPATH $$CMDPATH crd paths=./src/github.com/alsenz/esl-games/pkg/k8s/api/v1alpha1/ output:crd:dir=$$OUTDIR",
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
    k8s_container_push(
        name = "games-server-image." + CONFIGURATION_NAME + ".push",
        image = "//cmd/game-server:game-server-image",
        image_digest_tag = True,
        registry = REGISTRY,
        stamp = False,
    )
    for CONFIGURATION_NAME, (_, _, _, REGISTRY) in CONFIGURATION_TABLE.items()
]

# A genrule that runs the image push for the games server and then uses the digest to generate an overlay for the manager with
# the new image digest so it can correctly spin up deployments using it.
[
    #TODO: come up with a solution that is less linux-specific here - use of head, tail, and sed!
    genrule(
        name = "lesson-controller-gameserver-image-overlay." + CONFIGURATION_NAME + ".generate",
        srcs = [
            "games-server-image." + CONFIGURATION_NAME + ".push.digest",
            "overlays/lesson-operator/manager.yaml.tpl",
        ],
        outs = ["overlays/lesson-operator/manager.yaml"],
        cmd = "HOME=/tmp $(execpath games-server-image." + CONFIGURATION_NAME + ".push);" +
              "DIGEST=$$(head -c 14 $(execpath games-server-image." + CONFIGURATION_NAME + ".push.digest) | tail -c 7);" +
              "FULL_IMAGE=" + REGISTRY + "\\\/cmd\\\/game-server\\\/game-server-image:$$DIGEST;" +
              "sed \"s/LATEST_IMAGE/$$FULL_IMAGE/g\" $(execpath overlays/lesson-operator/manager.yaml.tpl) > $@",
        tools = ["games-server-image." + CONFIGURATION_NAME + ".push"],
    )
    for CONFIGURATION_NAME, (_, _, _, REGISTRY) in CONFIGURATION_TABLE.items()
]

[
    k8s_deploy(
        name = "lesson-controller." + DEP_MODE + "." + CONFIGURATION_NAME,
        cluster = CLUSTER,
        configurations = [
            #TODO I have no idea if this is working, need to validate this!
            "manifests/lesson-operator/kustomization.yaml",
            "manifests/lesson-operator/lesson_controller_manager_config.yaml",
        ],
        gitops = False,
        image_digest_tag = True,
        image_registry = REGISTRY,
        images = {
            "lesson-controller": "//cmd/lesson-operator:lesson-operator-image",
        },
        manifests = ["manifests/lesson-operator/manager.yaml"],
        namespace = PRIMARY_NAMESPACE,
        objects = [":namespace.deps." + CONFIGURATION_NAME] if DEP_MODES[DEP_MODE] else [],
        patches = ["overlays/lesson-operator/manager.yaml"],
        user = USER,
        deps = ["games-server-image." + CONFIGURATION_NAME + ".push"] if DEP_MODES[DEP_MODE] else [],
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

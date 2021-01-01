load("@com_adobe_rules_gitops//gitops:defs.bzl", "k8s_deploy")

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

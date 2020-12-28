load("@com_adobe_rules_gitops//gitops:defs.bzl", "k8s_deploy")

#TODO create a nice big configuration table that we can apply and add overlays too, then we loop over each one for the rules below

CONFIGURATION_NAME = "microk8s"
PRIMARY_NAMESPACE = "esl-games"
USER = "admin"
CLUSTER = "microk8s-cluster"
MICROK8S_REGISTRY = "localhost:32000"
USER = "admin"

k8s_deploy(
    name = "namespace." + CONFIGURATION_NAME,
    cluster = CLUSTER,
    image_registry = MICROK8S_REGISTRY,
    manifests = ["manifests/namespace.yaml"],
    namespace = PRIMARY_NAMESPACE,
    user = USER
)

k8s_deploy(
    name = "postgres." + CONFIGURATION_NAME,
    objects = [":namespace." + CONFIGURATION_NAME],
    cluster = CLUSTER,
    image_registry = MICROK8S_REGISTRY,
    manifests = ["manifests/postgres.yaml"],
    namespace = PRIMARY_NAMESPACE,
    user = USER
)

k8s_deploy(
    name = "matcher." + CONFIGURATION_NAME,
    objects = [":namespace." + CONFIGURATION_NAME],
    cluster = CLUSTER,
    image_registry = MICROK8S_REGISTRY,
    manifests = ["manifests/matcher.yaml"],
    namespace = PRIMARY_NAMESPACE,
    user = USER,
    images = {
        "matcher-image": "@webapp//:webapp-image",
    },
)

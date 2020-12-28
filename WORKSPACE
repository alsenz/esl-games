


# Gitops setup
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

rules_gitops_version = "8d9416a36904c537da550c95dc7211406b431db9"

http_archive(
    name = "com_adobe_rules_gitops",
    sha256 = "25601ed932bab631e7004731cf81a40bd00c9a34b87c7de35f6bc905c37ef30d",
    strip_prefix = "rules_gitops-%s" % rules_gitops_version,
    urls = ["https://github.com/adobe/rules_gitops/archive/%s.zip" % rules_gitops_version],
)

load("@com_adobe_rules_gitops//gitops:deps.bzl", "rules_gitops_dependencies")

rules_gitops_dependencies()

load("@com_adobe_rules_gitops//gitops:repositories.bzl", "rules_gitops_repositories")

rules_gitops_repositories()




# Same workspace rules as in the webapp repo

# Bazel Javascript / Nodejs setup

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "84b1d11b1f3bda68c24d992dc6e830bca9db8fa12276f2ca7fcb7761c893976b",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/3.0.0-rc.1/rules_nodejs-3.0.0-rc.1.tar.gz"],
)

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")


# NOTE: this rule installs nodejs, npm, and yarn, but does NOT install
# your npm dependencies into your node_modules folder.
# You must still run the package manager to do this.
node_repositories(package_json = ["@webapp//:package.json"])

yarn_install(
    # Name this npm so that Bazel Label references look like @npm//package
    name = "npm",
    package_json = "@webapp//:package.json",
    yarn_lock = "@webapp//:yarn.lock",
)


# And now we add the webapp repo and ensure that the yarn rule runs on patch
webapp_version = "master"
http_archive(
    name = "webapp",
    strip_prefix = "esl-games-webapp-%s" % webapp_version,
    urls = ["https://github.com/alsenz/esl-games-webapp/archive/%s.zip" % webapp_version],
    patch_cmds = ["bazel run @nodejs//:yarn"],
)

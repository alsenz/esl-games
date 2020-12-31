#!/usr/bin/env bash


# Licenses: much of this is borrowed from bazel-bin/gazelle-runner.bash, which is Copyright 2017 The Bazel Authors. These snippets are licensed under the Apache License 2.0
# ... see the notice in bazel-bin/gazelle-runner.bash for a full copyright notice and License statement.

GOTOOL='external/go_sdk/bin/go'

# find_runfile prints the location of a runfile in the source workspace,
# either by reading the symbolic link or reading the runfiles manifest.
function find_runfile {
  local runfile=$1
  if [ -f "$runfile" ]; then
    readlink "$runfile"
    return
  fi
  runfile=$(echo "$runfile" | sed -e 's!^\(\.\./\|external/\)!!')
  if grep -q "^$runfile" MANIFEST; then
    grep "^$runfile" MANIFEST | head -n 1 | cut -d' ' -f2
    return
  fi
  # printing nothing indicates failure
}

gotool=$(find_runfile "$GOTOOL")
if [ -z "$gotool" ]; then
  echo "$0: warning: could not locate GOROOT used by rules_go" >&2
  return
fi

if [ -z "$1" ]; then
  >&2 echo "No arguments provided, in this configuration a package must be manually specified... (and any other optional arguments)";
  exit 1;
fi

env GIT_TERMINAL_PROMPT=1 $gotool get -u $@

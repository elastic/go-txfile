#!/bin/bash

set -euo pipefail

WORKSPACE="$(pwd)/bin"
HW_TYPE="$(uname -m)"
PLATFORM_TYPE="$(uname)"

with_go() {
    local go_version="${1}"
    echo "Setting up the Go environment..."
    create_workspace
    check_platform_architeture
    retry 5 curl -sL -o ${WORKSPACE}/gvm "https://github.com/andrewkroh/gvm/releases/download/${SETUP_GVM_VERSION}/gvm-${PLATFORM_TYPE}-${arch_type}"
    export PATH="${PATH}:${WORKSPACE}"
    chmod +x ${WORKSPACE}/gvm
    eval "$(gvm "$go_version")"
    go version
    which go
    export PATH="${PATH}:$(go env GOPATH):$(go env GOPATH)/bin"
}

install_go_dependencies() {
    local install_packages=(
            "github.com/magefile/mage@v1.15.0"
            "github.com/elastic/go-licenser@v0.4.1"
            "golang.org/x/tools/cmd/goimports@v0.14.0"
            "github.com/jstemmer/go-junit-report/v2@v2.1.0"
            "github.com/tebeka/go2xunit@v1.4.10"
    )
    create_workspace
    for pkg in "${install_packages[@]}"; do
        go install "${pkg}"
    done
}

create_workspace() {
    if [[ ! -d "${WORKSPACE}" ]]; then
    mkdir -p ${WORKSPACE}
    fi
}

check_platform_architeture() {
# for downloading the GVM and Terraform packages
  case "${HW_TYPE}" in
   "x86_64")
        arch_type="amd64"
        ;;
    "aarch64")
        arch_type="arm64"
        ;;
    "arm64")
        arch_type="arm64"
        ;;
    *)
    echo "The current platform/OS type is unsupported yet"
    ;;
  esac
}

retry() {
    local retries=$1
    shift
    local count=0
    until "$@"; do
        exit=$?
        wait=$((2 ** count))
        count=$((count + 1))
        if [ $count -lt "$retries" ]; then
            >&2 echo "Retry $count/$retries exited $exit, retrying in $wait seconds..."
            sleep $wait
        else
            >&2 echo "Retry $count/$retries exited $exit, no more retries left."
            return $exit
        fi
    done
    return 0
}

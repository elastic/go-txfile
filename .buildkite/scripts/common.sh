#!/bin/bash

set -euo pipefail

WORKSPACE="$(pwd)/bin"
HW_TYPE="$(uname -m)"

with_go() {
    echo "Setting up the Go environment..."
    create_workspace
    check_platform_architeture
    retry 5 curl -sL -o ${WORKSPACE}/gvm "https://github.com/andrewkroh/gvm/releases/download/${SETUP_GVM_VERSION}/gvm-${platform_type_lowercase}-${arch_type}"
    chmod +x ${WORKSPACE}/gvm
    eval "$(gvm $(cat .go-version))"
    go version
    which go
    export PATH="${PATH}:$(go env GOPATH):$(go env GOPATH)/bin"
}

install_go() {
    local go_version="${1:-latest}"
    local install_packages=(
            "github.com/magefile/mage"
            "github.com/elastic/go-licenser"
            "golang.org/x/tools/cmd/goimports"
            "github.com/jstemmer/go-junit-report"
            "gotest.tools/gotestsum"
    )
    create_workspace
    for pkg in "${install_packages[@]}"; do
        go install "${pkg}@${go_version}"
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
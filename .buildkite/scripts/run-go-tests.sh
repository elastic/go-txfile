#!/bin/bash

set -euxo pipefail

source .buildkite/scripts/common.sh

GO_VERSION=$1
PLATFORM=$2
OUT_FILE="build/test-report.out"

with_go $GO_VERSION
install_go_dependencies $GO_VERSION $PLATFORM

# Run the tests
set +e
mkdir -p build
mage -v test | tee ${OUT_FILE}
status=$?
go get -v -u github.com/tebeka/go2xunit
go2xunit -fail -input ${OUT_FILE} -output "build/junit-${GO_VERSION}.xml"

exit ${status}

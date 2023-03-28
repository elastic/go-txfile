#!/usr/bin/env bash
set -euxo pipefail

# Run the tests
set +e
export OUT_FILE="build/test-report.out"
mkdir -p build
mage -v test | tee ${OUT_FILE}
status=$?
go get -v -u github.com/tebeka/go2xunit
go2xunit -fail -input ${OUT_FILE} -output "build/junit-${GO_VERSION}-${RUNNER_OS}.xml"

exit ${status}

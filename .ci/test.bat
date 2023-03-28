mkdir -p build
SET OUT_FILE=build\output-report.out
mage -v test > %OUT_FILE% | type %OUT_FILE%
go get -v -u github.com/tebeka/go2xunit
go2xunit -fail -input %OUT_FILE% -output build\junit-%GO_VERSION%-%RUNNER_OS%.xml

#mkdir -p build
#SET OUT_FILE=build\output-report.out
#mage -v test > %OUT_FILE% | type %OUT_FILE%
#go get -v -u github.com/tebeka/go2xunit
#go2xunit -fail -input %OUT_FILE% -output build\junit-%GO_VERSION%.xml


# TODO set $PATH to point to $GOROOT/bin

$GoVersion=$args[0]
$OutFile="build\output-report.out"
$ErrorActionPreference = "Stop" # set -e

# Forcing to checkout again all the files with a correct autocrlf.
# Doing this here because we cannot set git clone options before.
function fixCRLF {
    Write-Host "-- Fixing CRLF in git checkout --"
    git config core.autocrlf input
    git rm --quiet --cached -r .
    git reset --quiet --hard
}

function withGolang($version) {
    Write-Host "-- Install golang --"
    choco install -y golang --version $version
    $choco = Convert-Path "$((Get-Command choco).Path)\..\.."
    Import-Module "$choco\helpers\chocolateyProfile.psm1"
    refreshenv
    go version
    go env
}

function installGoDependencies($version) {
    $installPackages = @(
        "github.com/magefile/mage"
        "github.com/elastic/go-licenser"
        "golang.org/x/tools/cmd/goimports"
        "github.com/jstemmer/go-junit-report/v2"
        "gotest.tools/gotestsum"
    )
    foreach ($pkg in $installPackages) {
        go install "$pkg@$version"
    }
}

fixCRLF
withGolang $env:GO_VERSION_CHOCO
installGoDependencies

$ErrorActionPreference = "Continue" # set +e

New-Item -ItemType Directory -Force -Path "build"
mage -v test | Out-File -FilePath $OutFile
go get -v -u github.com/tebeka/go2xunit
go2xunit -fail -input $OutFile -output build\junit-$GoVersion.xml

#gotestsum --format testname --junitfile junit-win-report.xml -- -v ./...

$EXITCODE=$LASTEXITCODE
$ErrorActionPreference = "Stop"

Exit $EXITCODE
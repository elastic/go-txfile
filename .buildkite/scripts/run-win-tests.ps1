# TODO set $PATH to point to $GOROOT/bin

$GoVersion=$args[0]
$OutFile="build\output-report.out"
$ErrorActionPreference = "Stop" # set -e

# Forcing to checkout again all the files with a correct autocrlf.
# Doing this here because we cannot set git clone options before.
function fixCRLF {
    Write-Host "-- Fixing CRLF in git checkout --"
    git config core.autocrlf true
    git rm --quiet --cached -r .
    git reset --quiet --hard
}

function withGolang($version) {
    Write-Host "-- Install golang $version --"
    choco install golang -y --version $version
    $choco = Convert-Path "$((Get-Command choco).Path)\..\.."
    Import-Module "$choco\helpers\chocolateyProfile.psm1"
    refreshenv
    go version
    go env

}

function installGoDependencies() {
    $installPackages = @(
        "github.com/magefile/mage@v1.15.0"
        "github.com/elastic/go-licenser@v0.4.1"
        "golang.org/x/tools/cmd/goimports@v0.14.0"
        "github.com/jstemmer/go-junit-report@v1.0.0"
        "github.com/tebeka/go2xunit@v1.4.10"
    )
    foreach ($pkg in $installPackages) {
        go install "$pkg"
    }
}

fixCRLF
withGolang $GoVersion
installGoDependencies

$ErrorActionPreference = "Continue" # set +e

New-Item -ItemType Directory -Force -Path "build"
#go get -v -u github.com/tebeka/go2xunit
mage test | go2xunit -fail -output build\junit-$GoVersion.xml

$EXITCODE=$LASTEXITCODE
$ErrorActionPreference = "Stop"

Exit $EXITCODE

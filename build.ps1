<#
.SYNOPSIS
PowerShell script for testing, building, and deploying the experiment.
#>

$PROJECT = "dataset"
$GIT_GROUP = "caltechlibrary"
$RELEASE_DATE = Get-Date -Format "yyyy-MM-dd"
$RELEASE_HASH = git log --pretty=format:'%h' -n 1
$BRANCH = git rev-parse --abbrev-ref HEAD

# Determine the version from codemeta.json
$VERSION = jq -r '.version' codemeta.json

$MAN_PAGES = @("dataset.1", "datasetd.1", "dsquery.1", "dsimporter.1")
$MAN_PAGES_MISC = @("datasetd_yaml.5", "datasetd_service.5", "datasetd_api.5")
$PROGRAMS = @("dataset", "datasetd", "dsquery", "dsimporter")
$PREFIX = $HOME

$EXT = if ($IsWindows) { ".exe" } else { "" }
$EXT_WEB = ".wasm"
$DIST_FOLDERS = "bin/*", "man/*"

function Invoke-Build {
    # Build version.go
    cmt codemeta.json version.go
    git add version.go

    # Build programs
    foreach ($program in $PROGRAMS) {
        New-Item -ItemType Directory -Path "bin" -Force | Out-Null
        go build -o "bin/$program$EXT" "cmd/$program/$program.go"
        & "./bin/$program" "-help" | Out-File "$program.1.md"
    }

    # Build man pages
    foreach ($manPage in $MAN_PAGES) {
        New-Item -ItemType Directory -Path "man/man1" -Force | Out-Null
        pandoc "$manPage.md" --from markdown --to man -s | Out-File "man/man1/$manPage"
    }

    foreach ($manPageMisc in $MAN_PAGES_MISC) {
        New-Item -ItemType Directory -Path "man/man5" -Force | Out-Null
        pandoc "$manPageMisc.md" --from markdown --to man -s | Out-File "man/man5/$manPageMisc"
    }

    # Build CITATION.cff and about.md
    cmt codemeta.json CITATION.cff
    cmt codemeta.json about.md

    # Build installer scripts
    cmt codemeta.json installer.sh
    #chmod 775 installer.sh
    git add -f installer.sh

    cmt codemeta.json installer.ps1
    #chmod 775 installer.ps1
    git add -f installer.ps1
}

function Invoke-Install {
    if (-not (Test-Path "$PREFIX/bin")) {
        New-Item -ItemType Directory -Path "$PREFIX/bin" -Force | Out-Null
    }

    Write-Output "Installing programs in $PREFIX/bin"
    foreach ($program in $PROGRAMS) {
        if (Test-Path "./bin/$program") {
            Move-Item -Path "./bin/$program" -Destination "$PREFIX/bin/$program" -Force -Verbose
        }
    }

    Write-Output "Make sure $PREFIX/bin is in your PATH"
    foreach ($manPage in $MAN_PAGES) {
        if (Test-Path "./man/man1/$manPage") {
            Copy-Item -Path "./man/man1/$manPage" -Destination "$PREFIX/man/man1/$manPage" -Force -Verbose
        }
    }

    Write-Output "Make sure $PREFIX/man is in your MANPATH"
}

function Invoke-Uninstall {
    Write-Output "Removing programs in $PREFIX/bin"
    foreach ($program in $PROGRAMS) {
        if (Test-Path "$PREFIX/bin/$program") {
            Remove-Item -Path "$PREFIX/bin/$program" -Force -Verbose
        }
    }

    Write-Output "Removing manpages in $PREFIX/man"
    foreach ($manPage in $MAN_PAGES) {
        if (Test-Path "$PREFIX/man/man1/$manPage") {
            Remove-Item -Path "$PREFIX/man/man1/$manPage" -Force -Verbose
        }
    }
}

function Invoke-Check {
    go vet *.go
    Get-ChildItem -Directory | ForEach-Object {
        if (Test-Path "$_/*.go") {
            Push-Location $_
            go vet *.go
            Pop-Location
        }
    }
}

function Invoke-Test {
    go test
}

function Invoke-Clean {
    go clean
    Remove-Item -Path "bin" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "dist" -Recurse -Force -ErrorAction SilentlyContinue
    Remove-Item -Path "testout" -Recurse -Force -ErrorAction SilentlyContinue
    Get-ChildItem -Directory | ForEach-Object {
        Remove-Item -Path "$_/testout" -Recurse -Force -ErrorAction SilentlyContinue
    }
    go clean -r
}

function Invoke-Dist {
    param($os, $arch, $suffix)
    $distPath = "dist/bin"
    New-Item -ItemType Directory -Path $distPath -Force | Out-Null

    foreach ($program in $PROGRAMS) {
        $env:GOOS = $os
        $env:GOARCH = $arch
        go build -o "$distPath/$program$EXT" "cmd/$program/*.go"
    }

    Push-Location dist
    zip -r "${PROJECT}-v${VERSION}-${suffix}.zip" LICENSE codemeta.json CITATION.cff *.md $DIST_FOLDERS
    Pop-Location

    Remove-Item -Path $distPath -Recurse -Force
}

function Invoke-DistributeDocs {
    Remove-Item -Path "dist" -Recurse -Force -ErrorAction SilentlyContinue
    New-Item -ItemType Directory -Path "dist" -Force | Out-Null

    Copy-Item -Path "codemeta.json", "CITATION.cff", "README.md", "LICENSE", "INSTALL.md", "installer.sh", "installer.ps1" -Destination "dist" -Verbose
    Copy-Item -Path "man" -Destination "dist" -Recurse -Verbose
}

function Invoke-Release {
    Invoke-Clean
    Invoke-Build
    Invoke-Dist "linux" "amd64" "Linux-x86_64"
    Invoke-Dist "linux" "arm64" "Linux-aarch64"
    Invoke-Dist "linux" "arm" "Linux-armv7l"
    Invoke-Dist "windows" "amd64" "Windows-x86_64"
    Invoke-Dist "windows" "arm64" "Windows-arm64"
    Invoke-Dist "darwin" "amd64" "macOS-x86_64"
    Invoke-Dist "darwin" "arm64" "macOS-arm64"
    .\release.ps1
}

function Invoke-Status {
    git status
}

function Invoke-Save {
    param([string]$msg = "Quick Save")
    git commit -am $msg
    git push origin $BRANCH
}

function Invoke-Publish {
    .\publish.ps1
}

function Invoke-LogHash {
    git log --pretty=format:'%h' -n 1
}

# Example usage:
# Invoke-Build
# Invoke-Install
# Invoke-Test
Invoke-Build

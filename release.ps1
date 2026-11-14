#!/usr/bin/env pwsh
#
# release.ps1 creates a draft GitHub release from dist/*.zip. It does what
# release.bash does, with the same checks, messages and exit statuses; the
# shared test driver (release_test.bash, RELEASE_IMPL=ps1) runs the same checks
# against both. Run `make release` first: it regenerates the version files,
# builds the zip files and checks them (release-check), so by the time this
# script runs the checks have already passed.
#
# Usage: ./release.ps1 [--dry-run] [--yes]
#
# The release tag and notes come from codemeta.json (version, releaseNotes).
# The script refuses to run if anything is uncommitted other than the files
# `make release` and the author's version bump change (the prep files below), or
# if commits are not pushed. After you confirm, it commits those files as
# "prep for vX.Y.Z", pushes, writes a checksums file next to the zip files and
# creates a DRAFT release carrying the zips and the checksums. Publish the draft
# from the repository's releases page.
#
# Options use the same spelling as release.bash (not PowerShell parameters), so
# a command line works with either script.

$ErrorActionPreference = 'Stop'

$E_NEGATIVE = 1; $E_USAGE = 2; $E_DATA = 65; $E_NO_INPUT = 66; $E_UNAVAILABLE = 69; $E_IO = 74

$Usage = @'
release.ps1 creates a draft GitHub release from dist/*.zip. Run `make
release` first: it regenerates the version files, builds the zip files and
checks them (release-check), so by the time this script runs the checks have
already passed.

Usage: ./release.ps1 [--dry-run] [--yes]

The release tag and notes come from codemeta.json (version, releaseNotes).
The script refuses to run if anything is uncommitted other than the files
`make release` and the author's version bump change (the prep files), or if
commits are not pushed. After you confirm, it commits those files as
"prep for vX.Y.Z", pushes, writes a checksums file next to the zip files and
creates a DRAFT release carrying the zips and the checksums. Publish the draft
from the repository's releases page.

Options:
  --dry-run   run every check and print what would happen; create nothing
  --yes       do not ask for confirmation
  -h, --help  this text

Needs: PowerShell 7, git, gh.

EXIT STATUS
  0   draft release created (or --dry-run checks passed)
  1   a check refused: uncommitted changes beyond the release files, unpushed
      commits, the tag points elsewhere, the release already exists, or the
      prompt was declined
  2   usage: unknown option or surplus argument
  65  codemeta.json has no usable version (missing or unexpected characters)
  66  an input is missing: codemeta.json, dist/*.zip, or a required command
  69  git could not push, or gh could not create the release
  74  git could not commit the release files
'@

# The files a release legitimately changes: the author's version bump and what
# `make release` regenerates. Anything else uncommitted stops the release.
$PrepFiles = @('codemeta.json', 'version.go', 'CITATION.cff', 'installer.sh', 'installer.ps1')

function Die([int]$Code, [string]$Message) {
    [Console]::Error.WriteLine("release.ps1: $Message")
    exit $Code
}

function Test-PrepFile([string]$Path) {
    if ($Path -cmatch '\Adocs/.*\..\.md\z') { return $true }
    return ($PrepFiles -ccontains $Path)
}

$DryRun = $false
$AssumeYes = $false
foreach ($arg in $args) {
    switch -CaseSensitive ($arg) {
        '--dry-run' { $DryRun = $true }
        '--yes' { $AssumeYes = $true }
        '-h' { Write-Output $Usage; exit 0 }
        '--help' { Write-Output $Usage; exit 0 }
        default {
            if ($arg.StartsWith('-')) { Die $E_USAGE "usage: unknown option $arg (see --help)" }
            Die $E_USAGE "usage: unexpected argument $arg (see --help)"
        }
    }
}

# Required commands, all reported at once.
$missing = ''
foreach ($cmd in 'gh', 'git') {
    if (-not (Get-Command $cmd -CommandType Application -ErrorAction SilentlyContinue)) { $missing += " $cmd" }
}
if ($missing -ne '') { Die $E_NO_INPUT "required command(s) not found:$missing" }

if (-not (Test-Path -LiteralPath 'codemeta.json')) { Die $E_NO_INPUT "codemeta.json not found in $((Get-Location).Path)" }
try {
    $meta = Get-Content -Raw -LiteralPath 'codemeta.json' | ConvertFrom-Json
} catch {
    Die $E_DATA 'codemeta.json is not valid JSON'
}
$version = "$($meta.version)"
if ($version -eq '') { Die $E_DATA 'codemeta.json has no version' }
if ($version -cnotmatch '\A[0-9a-zA-Z._-]+\z') { Die $E_DATA "version contains unexpected characters: $version" }
$tag = "v$version"
$name = "$($meta.name)"
if ($name -eq '') { $name = Split-Path -Leaf (Get-Location).Path }

$zips = @(Get-ChildItem -Path 'dist' -Filter '*.zip' -File -ErrorAction SilentlyContinue | Sort-Object Name | ForEach-Object { "dist/$($_.Name)" })
if ($zips.Count -eq 0) { Die $E_NO_INPUT 'no dist/*.zip; run make release first' }

# A release is built from what is on GitHub: clean tree, everything pushed.
git rev-parse --is-inside-work-tree *> $null
if ($LASTEXITCODE -ne 0) { Die $E_NO_INPUT 'not inside a git work tree' }
$dirtyPrep = @()
foreach ($line in @(git status --porcelain)) {
    if ([string]::IsNullOrEmpty($line)) { continue }
    $changed = $line.Substring(3)
    # Only a plain modification of a release file is allowed (" M" or "M ").
    if ((@(' M', 'M ', 'MM') -ccontains $line.Substring(0, 2)) -and (Test-PrepFile $changed)) {
        $dirtyPrep += $changed
        continue
    }
    Die $E_NEGATIVE "uncommitted changes beyond the release files ($changed); commit or revert them first"
}
$upstream = git rev-parse --abbrev-ref --symbolic-full-name '@{u}' 2>$null
if ($LASTEXITCODE -ne 0) { Die $E_NEGATIVE 'commits not pushed: the branch has no upstream' }
$ahead = git rev-list --count "${upstream}..HEAD"
if ("$ahead" -ne '0') { Die $E_NEGATIVE "commits not pushed to $upstream; push first" }

$headHash = (git rev-parse HEAD)
$tagHash = git rev-parse -q --verify "refs/tags/${tag}^{commit}" 2>$null
if ($LASTEXITCODE -eq 0 -and "$tagHash" -ne "$headHash") {
    Die $E_NEGATIVE "tag $tag exists and points at $("$tagHash".Substring(0, 7)), not HEAD $($headHash.Substring(0, 7))"
}
& gh release view $tag *> $null
if ($LASTEXITCODE -eq 0) { Die $E_NEGATIVE "release $tag already exists on GitHub" }

$checksums = "$name-$tag-checksums.txt"
$notes = "$($meta.releaseNotes)"
Write-Output "tag:       $tag (commit $($headHash.Substring(0, 7)))"
Write-Output 'release:   draft'
Write-Output 'files:'
foreach ($z in $zips) { Write-Output "  $z" }
Write-Output "  dist/$checksums (to be written)"
Write-Output "notes:     $($notes.Substring(0, [Math]::Min(200, $notes.Length)))"

if ($dirtyPrep.Count -gt 0) { Write-Output "would commit and push: $($dirtyPrep -join ' ')" }

if ($DryRun) {
    Write-Output 'dry run: nothing created'
    exit 0
}

if (-not $AssumeYes) {
    [Console]::Out.Write("Create the draft release $tag on GitHub? (y/N) ")
    $answer = [Console]::In.ReadLine()
    if ($answer -cne 'y') { Die $E_NEGATIVE 'release cancelled; nothing sent to GitHub' }
}

if ($dirtyPrep.Count -gt 0) {
    git add -- @dirtyPrep
    git commit -q -m "prep for $tag"
    if ($LASTEXITCODE -ne 0) { Die $E_IO 'git could not commit the release files' }
    git push -q
    if ($LASTEXITCODE -ne 0) { Die $E_UNAVAILABLE 'git push failed; the release files are committed locally, nothing sent to GitHub' }
    Write-Output "committed and pushed: $($dirtyPrep -join ' ')"
}

$sums = foreach ($z in $zips) {
    $hash = (Get-FileHash -Algorithm SHA256 -LiteralPath $z).Hash.ToLowerInvariant()
    "$hash  $(Split-Path -Leaf $z)"
}
[System.IO.File]::WriteAllText((Join-Path 'dist' $checksums), (($sums -join "`n") + "`n"))
Write-Output "checksums written to dist/$checksums"

$notesFile = [System.IO.Path]::GetTempFileName()
try {
    [System.IO.File]::WriteAllText($notesFile, $notes + "`n")
    & gh release create $tag --draft --notes-file $notesFile @zips "dist/$checksums"
    if ($LASTEXITCODE -ne 0) { Die $E_UNAVAILABLE "gh could not create the release $tag" }
} finally {
    Remove-Item -LiteralPath $notesFile -Force -ErrorAction SilentlyContinue
}

Write-Output "draft release $tag created; review and publish it on the repository's releases page"

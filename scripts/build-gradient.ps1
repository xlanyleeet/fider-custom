param([string]$Image = 'gradient-fider:local')
$ErrorActionPreference = 'Stop'
$repoPath = Split-Path -Parent $PSScriptRoot
$commitHash = git -C $repoPath rev-parse HEAD
if ($LASTEXITCODE -ne 0) { throw 'Cannot read the source commit.' }
$version = 'gradient-' + $commitHash.Substring(0, 12)
docker build --platform linux/amd64 --build-arg "COMMITHASH=$commitHash" --build-arg "VERSION=$version" --tag $Image $repoPath
if ($LASTEXITCODE -ne 0) { throw 'Docker build failed.' }

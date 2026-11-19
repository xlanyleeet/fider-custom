param([switch]$SkipBuild)
$ErrorActionPreference = 'Stop'
$repoPath = Split-Path -Parent $PSScriptRoot
$envPath = Join-Path $repoPath 'deploy/.env.preview'

function New-RandomHex([int]$Length) {
    $buffer = New-Object byte[] $Length
    $rng = [Security.Cryptography.RandomNumberGenerator]::Create()
    try { $rng.GetBytes($buffer) } finally { $rng.Dispose() }
    return ([BitConverter]::ToString($buffer)).Replace('-', '').ToLowerInvariant()
}

if (-not (Test-Path -LiteralPath $envPath)) {
    $content = [IO.File]::ReadAllText((Join-Path $repoPath 'deploy/.env.example'))
    $content = $content.Replace('GENERATE_POSTGRES_PASSWORD', (New-RandomHex 32))
    $content = $content.Replace('GENERATE_JWT_SECRET', (New-RandomHex 64))
    [IO.File]::WriteAllText($envPath, $content, (New-Object Text.UTF8Encoding $false))
}

$env:FIDER_COMMIT = git -C $repoPath rev-parse HEAD
if ($LASTEXITCODE -ne 0) { throw 'Cannot read the source commit.' }
$env:FIDER_VERSION = 'gradient-' + $env:FIDER_COMMIT.Substring(0, 12)
$composeArgs = @('compose', '--env-file', $envPath, '-f', (Join-Path $repoPath 'deploy/compose.preview.yaml'), 'up', '-d', '--wait', '--wait-timeout', '180')
if (-not $SkipBuild) { $composeArgs += '--build' }
docker @composeArgs
if ($LASTEXITCODE -ne 0) { throw 'Preview failed to start.' }
Write-Host 'Preview: http://localhost:3001 (or HTTP_PORT from deploy/.env.preview)'
Write-Host 'Test inbox: http://localhost:8026 (or MAIL_UI_PORT from deploy/.env.preview)'

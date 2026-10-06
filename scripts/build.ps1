# Usage: .\scripts\build.ps1 all | windows | linux | ...
param(
    [Parameter(Position = 0)]
    [ValidateSet('help', 'all', 'local', 'windows', 'linux', 'linux-arm64', 'test', 'clean')]
    [string] $Target = 'help'
)

$ErrorActionPreference = 'Stop'
$Root = Split-Path -Parent (Split-Path -Parent $MyInvocation.MyCommand.Path)
Set-Location $Root

$Dist = if ($env:DIST) { $env:DIST } else { 'dist' }
$LdFlags = '-s -w'

function Show-Help {
    @"
Usage: .\scripts\build.ps1 <target>

Targets:
  all           Windows + Linux amd64
  local         binary for this machine
  windows       dist\lexigen-windows-amd64.exe
  linux         dist\lexigen-linux-amd64
  linux-arm64   dist\lexigen-linux-arm64
  test          go test ./...
  clean         remove dist\
"@
}

function Go-Build {
    param([string] $GoOS, [string] $GoArch, [string] $Out)
    $env:GOOS = $GoOS
    $env:GOARCH = $GoArch
    $env:CGO_ENABLED = '0'
    go build -ldflags $LdFlags -o $Out .
    Remove-Item Env:GOOS, Env:GOARCH, Env:CGO_ENABLED -ErrorAction SilentlyContinue
}

New-Item -ItemType Directory -Force -Path $Dist | Out-Null

switch ($Target) {
    'help' { Show-Help }
    'all' {
        Go-Build 'windows' 'amd64' (Join-Path $Dist 'lexigen-windows-amd64.exe')
        Go-Build 'linux' 'amd64' (Join-Path $Dist 'lexigen-linux-amd64')
        Write-Host "built: $Dist\lexigen-windows-amd64.exe $Dist\lexigen-linux-amd64"
    }
    'local' {
        $out = Join-Path $Dist 'lexigen.exe'
        go build -ldflags $LdFlags -o $out .
        Write-Host "built: $out"
    }
    'windows' {
        $out = Join-Path $Dist 'lexigen-windows-amd64.exe'
        Go-Build 'windows' 'amd64' $out
        Write-Host "built: $out"
    }
    'linux' {
        $out = Join-Path $Dist 'lexigen-linux-amd64'
        Go-Build 'linux' 'amd64' $out
        Write-Host "built: $out"
    }
    'linux-arm64' {
        $out = Join-Path $Dist 'lexigen-linux-arm64'
        Go-Build 'linux' 'arm64' $out
        Write-Host "built: $out"
    }
    'test' { go test ./... }
    'clean' {
        Remove-Item -Recurse -Force $Dist -ErrorAction SilentlyContinue
        Write-Host "removed $Dist\"
    }
}

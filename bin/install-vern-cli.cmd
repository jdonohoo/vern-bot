@echo off
REM install-vern-cli: Auto-download the vern CLI binary for Windows
REM
REM Downloads from GitHub releases and installs to {plugin_root}\go\bin\vern.exe
REM Requires PowerShell (available on all Windows 10+).

setlocal enabledelayedexpansion

set "REPO=jdonohoo/vern-bot"
set "SCRIPT_DIR=%~dp0"
set "INSTALL_DIR=%SCRIPT_DIR%..\go\bin"
set "BINARY_NAME=vern.exe"

REM Detect architecture
set "ARCH=amd64"
if "%PROCESSOR_ARCHITECTURE%"=="ARM64" set "ARCH=arm64"

set "ARCHIVE_NAME=vern_windows_%ARCH%.zip"

echo [vern-cli] Detecting platform: windows/%ARCH%
echo [vern-cli] Fetching latest release info...

REM Use PowerShell to fetch release info and download
powershell -NoProfile -ExecutionPolicy Bypass -Command ^
  "$ErrorActionPreference = 'Stop';" ^
  "$release = Invoke-RestMethod -Uri 'https://api.github.com/repos/%REPO%/releases/latest';" ^
  "$asset = $release.assets | Where-Object { $_.name -eq '%ARCHIVE_NAME%' };" ^
  "if (-not $asset) { Write-Error 'No release found for %ARCHIVE_NAME%'; exit 1 };" ^
  "$url = $asset.browser_download_url;" ^
  "Write-Host '[vern-cli] Downloading:' $url;" ^
  "$tmp = Join-Path $env:TEMP 'vern-cli-download';" ^
  "if (Test-Path $tmp) { Remove-Item $tmp -Recurse -Force };" ^
  "New-Item -ItemType Directory -Path $tmp -Force | Out-Null;" ^
  "$zipPath = Join-Path $tmp '%ARCHIVE_NAME%';" ^
  "Invoke-WebRequest -Uri $url -OutFile $zipPath;" ^
  "Write-Host '[vern-cli] Extracting...';" ^
  "Expand-Archive -Path $zipPath -DestinationPath $tmp -Force;" ^
  "$binary = Get-ChildItem -Path $tmp -Filter '%BINARY_NAME%' -Recurse | Select-Object -First 1;" ^
  "if (-not $binary) { Write-Error 'Binary not found in archive'; exit 1 };" ^
  "$installDir = '%INSTALL_DIR%';" ^
  "if (-not (Test-Path $installDir)) { New-Item -ItemType Directory -Path $installDir -Force | Out-Null };" ^
  "$dest = Join-Path $installDir '%BINARY_NAME%';" ^
  "Copy-Item $binary.FullName $dest -Force;" ^
  "Write-Host '[vern-cli] Installed:' $dest;" ^
  "Remove-Item $tmp -Recurse -Force;" ^
  "& $dest --version"

if errorlevel 1 (
    echo [vern-cli] Error: Download failed. You may need to build from source: cd go ^&^& go build -o bin\vern.exe .\cmd\vern 1>&2
    exit /b 1
)

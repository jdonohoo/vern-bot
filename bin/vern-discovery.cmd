@echo off
REM vern-discovery: Run the full Vern Discovery Pipeline (Windows)
REM
REM Delegates to the Go vern CLI binary. Auto-downloads if not present.
REM
REM Usage: vern-discovery [flags] "<idea>" [discovery_dir]

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "VERN_CLI=%SCRIPT_DIR%..\go\bin\vern.exe"

REM Auto-build if Go binary doesn't exist
if not exist "%VERN_CLI%" (
    where go >nul 2>nul
    if not errorlevel 1 (
        if exist "%SCRIPT_DIR%..\go\cmd\vern\main.go" (
            echo [vern-discovery] Building vern CLI from source... 1>&2
            pushd "%SCRIPT_DIR%..\go"
            go build -o bin\vern.exe .\cmd\vern 2>nul
            popd
        )
    )
)

REM Auto-download if still not available
if not exist "%VERN_CLI%" (
    if exist "%SCRIPT_DIR%install-vern-cli.cmd" (
        call "%SCRIPT_DIR%install-vern-cli.cmd"
    )
)

if not exist "%VERN_CLI%" (
    echo [vern-discovery] Error: vern CLI not available. Run: cd go ^&^& go build -o bin\vern.exe .\cmd\vern 1>&2
    exit /b 1
)

"%VERN_CLI%" discovery %*

@echo off
REM vernhole: Summon random Vern personas for chaotic discovery (Windows)
REM
REM Delegates to the Go vern CLI binary. Auto-downloads if not present.
REM
REM Usage: vernhole [flags] "<idea>"

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "VERN_CLI=%SCRIPT_DIR%..\go\bin\vern.exe"

REM Auto-build if Go binary doesn't exist
if not exist "%VERN_CLI%" (
    where go >nul 2>nul
    if not errorlevel 1 (
        if exist "%SCRIPT_DIR%..\go\cmd\vern\main.go" (
            echo [vernhole] Building vern CLI from source... 1>&2
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
    echo [vernhole] Error: vern CLI not available. Run: cd go ^&^& go build -o bin\vern.exe .\cmd\vern 1>&2
    exit /b 1
)

"%VERN_CLI%" hole %*

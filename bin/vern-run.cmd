@echo off
REM vern-run: Spawn different LLM subprocesses for Vern-Bot (Windows)
REM
REM Delegates to the Go vern CLI binary. Auto-downloads if not present.
REM
REM Usage: vern-run <llm> "<prompt>" [output_file] [persona]

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "VERN_CLI=%SCRIPT_DIR%..\go\bin\vern.exe"

REM Auto-build if Go binary doesn't exist
if not exist "%VERN_CLI%" (
    where go >nul 2>nul
    if not errorlevel 1 (
        if exist "%SCRIPT_DIR%..\go\cmd\vern\main.go" (
            echo [vern-run] Building vern CLI from source... 1>&2
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
    echo [vern-run] Error: vern CLI not available. Run: cd go ^&^& go build -o bin\vern.exe .\cmd\vern 1>&2
    exit /b 1
)

REM Map positional args to Go CLI flags
set "LLM=%~1"
set "PROMPT=%~2"
set "OUTPUT_FILE=%~3"
set "PERSONA=%~4"

if "%LLM%"=="" (
    echo Usage: vern-run ^<llm^> "^<prompt^>" [output_file] [persona]
    echo   llm: claude ^| codex ^| gemini
    exit /b 1
)
if "%PROMPT%"=="" (
    echo Usage: vern-run ^<llm^> "^<prompt^>" [output_file] [persona]
    exit /b 1
)

set "ARGS=run %LLM% "%PROMPT%""
if not "%OUTPUT_FILE%"=="" set "ARGS=%ARGS% --output "%OUTPUT_FILE%""
if not "%PERSONA%"=="" set "ARGS=%ARGS% --persona %PERSONA%"

"%VERN_CLI%" %ARGS%

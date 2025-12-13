@echo off
REM Debug build script for ClaudeCompanion
REM Builds application with console window for debugging

echo ==========================================
echo ClaudeCompanion - Debug Build Script
echo ==========================================
echo.

REM Navigate to project root
cd /d "%~dp0.."

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go 1.21 or higher from https://golang.org/dl/
    pause
    exit /b 1
)

echo [1/3] Installing dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)

echo.
echo [2/3] Building application (debug mode with console)...
go build -o dist/claudecompanion-debug.exe ./cmd/claudecompanion
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to build application
    pause
    exit /b 1
)

echo.
echo [3/3] Copying required files...
copy icon.ico dist\ >nul 2>nul
if not exist dist\config.yaml (
    copy config.yaml.example dist\config.yaml >nul 2>nul
    echo Created dist\config.yaml from example
) else (
    echo Config file already exists, skipping
)

echo.
echo ==========================================
echo Debug build completed successfully!
echo ==========================================
echo.
echo Output: dist\claudecompanion-debug.exe
echo.
echo This version shows console window with logs.
echo Run it to see real-time debug output.
echo.

REM Ask if user wants to run it now
set /p RUN="Run debug version now? (Y/N): "
if /i "%RUN%"=="Y" (
    echo.
    echo Starting claudecompanion-debug.exe...
    echo Press Ctrl+C to stop
    echo.
    dist\claudecompanion-debug.exe
) else (
    echo.
    echo You can run it manually: dist\claudecompanion-debug.exe
    pause
)

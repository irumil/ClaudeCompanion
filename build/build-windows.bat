@echo off
REM Build script for Windows

echo Building ClaudeCompanion for Windows...

REM Change to project root
cd ..

REM Set variables
set OUTPUT_DIR=dist
set APP_NAME=claudecompanion.exe
set MAIN_PACKAGE=./cmd/claudecompanion

REM Create output directory
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

REM Download dependencies
echo Downloading dependencies...
go mod tidy
go mod download

REM Build for Windows
echo Compiling for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

go build -ldflags="-s -w -H windowsgui" -o %OUTPUT_DIR%\%APP_NAME% %MAIN_PACKAGE%

if %ERRORLEVEL% EQU 0 (
    echo.
    echo ========================================
    echo Build successful!
    echo Output: %OUTPUT_DIR%\%APP_NAME%
    echo ========================================
    echo.
    echo To run the application:
    echo   %OUTPUT_DIR%\%APP_NAME%
    echo.
    echo Config file will be created at:
    echo   %%APPDATA%%\ClaudeCompanion\config.yaml
) else (
    echo.
    echo ========================================
    echo Build failed!
    echo ========================================
    exit /b 1
)

pause

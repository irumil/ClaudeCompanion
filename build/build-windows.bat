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

REM Prepare icon
echo Preparing icon for embedding...
copy icon.ico cmd\claudecompanion\icon.ico >nul 2>nul

REM Generate Windows resource file
echo Generating Windows resource file...
cd cmd\claudecompanion
"%USERPROFILE%\go\bin\rsrc.exe" -ico "..\..\icon.ico" -o rsrc_windows_amd64.syso 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: rsrc not found or failed, icon may not be embedded in exe
)
cd ..\..

REM Build for Windows
echo Compiling for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0

go build -ldflags="-s -w -H windowsgui" -o %OUTPUT_DIR%\%APP_NAME% %MAIN_PACKAGE%

if %ERRORLEVEL% EQU 0 (
    REM Copy required files
    echo Copying required files...
    if not exist %OUTPUT_DIR%\config.yaml (
        copy config.yaml.example %OUTPUT_DIR%\config.yaml >nul 2>nul
    )

    echo.
    echo ========================================
    echo Build successful!
    echo ========================================
    echo.
    echo Output files:
    echo   %OUTPUT_DIR%\%APP_NAME% ^(with embedded icon^)
    echo   %OUTPUT_DIR%\config.yaml ^(configuration^)
    echo.
    echo To run the application:
    echo   %OUTPUT_DIR%\%APP_NAME%
) else (
    echo.
    echo ========================================
    echo Build failed!
    echo ========================================
    exit /b 1
)

pause

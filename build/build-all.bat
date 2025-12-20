@echo off
REM Full build script for ClaudeCompanion
REM Builds both the application and browser extension

echo ==========================================
echo ClaudeCompanion - Full Build Script
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

echo [1/6] Installing dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)

echo.
echo [2/6] Preparing icon for embedding...
REM Copy icon to cmd/claudecompanion for go:embed
copy icon.ico cmd\claudecompanion\icon.ico >nul 2>nul

REM Generate Windows resource file for exe icon
cd cmd\claudecompanion
"%USERPROFILE%\go\bin\rsrc.exe" -ico "..\..\icon.ico" -o rsrc_windows_amd64.syso
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: Failed to generate Windows resource, continuing anyway...
)
cd ..\..

echo.
echo [3/6] Building application (release mode)...
go build -ldflags "-H windowsgui" -o dist/claudecompanion.exe ./cmd/claudecompanion
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to build application
    pause
    exit /b 1
)

echo.
echo [4/6] Copying required files...
if not exist dist\config.yaml (
    copy config.yaml.example dist\config.yaml >nul 2>nul
    echo Created dist\config.yaml from example
) else (
    echo Config file already exists, skipping
)

echo.
echo [5/6] Building browser extension...
cd build
call package-extension.bat
cd ..

echo.
echo [6/6] Verifying build...
if exist dist\claudecompanion.exe (
    echo [OK] claudecompanion.exe (with embedded icon)
) else (
    echo [FAIL] claudecompanion.exe missing
)

if exist dist\claudecompanion-extension.zip (
    echo [OK] claudecompanion-extension.zip
) else (
    echo [FAIL] claudecompanion-extension.zip missing
)

if exist dist\config.yaml (
    echo [OK] config.yaml
) else (
    echo [FAIL] config.yaml missing
)

echo.
echo ==========================================
echo Build completed successfully!
echo ==========================================
echo.
echo Output files in 'dist' folder:
echo   - claudecompanion.exe (Desktop application with embedded icon)
echo   - claudecompanion-extension.zip (Browser extension)
echo   - config.yaml (Configuration file)
echo.
echo Next steps:
echo   1. Configure dist\config.yaml (optional)
echo   2. Run dist\claudecompanion.exe
echo   3. Install extension from dist\claudecompanion-extension.zip
echo.
pause

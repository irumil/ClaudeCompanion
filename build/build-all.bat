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

echo [1/5] Installing dependencies...
go mod download
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)

echo.
echo [2/5] Building application (release mode)...
go build -ldflags "-H windowsgui" -o dist/claudecompanion.exe ./cmd/claudecompanion
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: Failed to build application
    pause
    exit /b 1
)

echo.
echo [3/5] Copying required files...
copy icon.ico dist\ >nul 2>nul
if not exist dist\config.yaml (
    copy config.yaml.example dist\config.yaml >nul 2>nul
    echo Created dist\config.yaml from example
) else (
    echo Config file already exists, skipping
)

echo.
echo [4/5] Building browser extension...
cd build
call package-extension.bat
cd ..

echo.
echo [5/5] Verifying build...
if exist dist\claudecompanion.exe (
    echo [OK] claudecompanion.exe
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

if exist dist\icon.ico (
    echo [OK] icon.ico
) else (
    echo [FAIL] icon.ico missing
)

echo.
echo ==========================================
echo Build completed successfully!
echo ==========================================
echo.
echo Output files in 'dist' folder:
echo   - claudecompanion.exe (Desktop application)
echo   - claudecompanion-extension.zip (Browser extension)
echo   - config.yaml (Configuration file)
echo   - icon.ico (Application icon)
echo.
echo Next steps:
echo   1. Configure dist\config.yaml (optional)
echo   2. Run dist\claudecompanion.exe
echo   3. Install extension from dist\claudecompanion-extension.zip
echo.
pause

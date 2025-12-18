@echo off
REM Package ClaudeCompanion Firefox Extension
echo ==========================================
echo Packaging ClaudeCompanion Firefox Extension
echo ==========================================
echo.

cd /d "%~dp0..\extension"

REM Extract version from manifest.json
for /f "tokens=2 delims=:, " %%a in ('findstr /C:"\"version\"" manifest.json') do set VERSION=%%a
set VERSION=%VERSION:"=%

echo Version: %VERSION%
echo.

REM Remove old packages if exist
if exist "..\dist\claudecompanion-extension.zip" del "..\dist\claudecompanion-extension.zip"
if exist "..\dist\claudecompanion-%VERSION%.xpi" del "..\dist\claudecompanion-%VERSION%.xpi"
if exist "..\dist\claudecompanion-extension.xpi" del "..\dist\claudecompanion-extension.xpi"

REM Create temporary ZIP package with ALL required files
echo [1/3] Creating package...
powershell -Command "Compress-Archive -Path manifest.json,background.js,options.html,options.js,icon48.png,icon96.png,README.md,PRIVACY.md -DestinationPath ../dist/temp-extension.zip -Force"

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ERROR: Failed to create package!
    pause
    exit /b 1
)

REM Rename ZIP to XPI
echo [2/3] Renaming to XPI...
move "..\dist\temp-extension.zip" "..\dist\claudecompanion-%VERSION%.xpi" >nul

REM Create copy without version for convenience
copy "..\dist\claudecompanion-%VERSION%.xpi" "..\dist\claudecompanion-extension.xpi" >nul

echo [3/3] Verifying package...
powershell -Command "$size = (Get-Item '../dist/claudecompanion-%VERSION%.xpi').Length; $kb = [math]::Round($size/1KB, 2); Write-Host \"Size: $kb KB\""

echo.
echo ==========================================
echo Extension packaged successfully!
echo ==========================================
echo.
echo Created files:
echo   - dist\claudecompanion-%VERSION%.xpi (versioned)
echo   - dist\claudecompanion-extension.xpi (latest)
echo.
echo Package includes:
echo   - Code: manifest.json, background.js, options.*
echo   - Icons: icon48.png, icon96.png
echo   - Docs: README.md, PRIVACY.md
echo.
echo Ready for upload to:
echo https://addons.mozilla.org/developers/
echo.

pause

@echo off
REM Package ClaudeCompanion Firefox Extension
echo Packaging ClaudeCompanion Firefox Extension...

cd /d "%~dp0..\extension"

REM Remove old package if exists
if exist "..\dist\claudecompanion-extension.zip" (
    del "..\dist\claudecompanion-extension.zip"
    echo Removed old package
)

REM Create new package
powershell -Command "Compress-Archive -Path manifest.json,background.js,options.html,options.js,icon48.png,icon96.png -DestinationPath ../dist/claudecompanion-extension.zip -Force"

if %ERRORLEVEL% EQU 0 (
    echo.
    echo Extension packaged successfully!
    echo Location: dist\claudecompanion-extension.zip
    echo.
    echo You can now upload this file to:
    echo https://addons.mozilla.org/developers/
) else (
    echo.
    echo ERROR: Failed to package extension!
    exit /b 1
)

pause

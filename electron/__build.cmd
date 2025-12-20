@echo off
pushd %~dp0
taskkill /f /im:go-boost-electron.exe 2>nul
if not exist package.json (
    if exist package.json.in (
        copy /y package.json.in package.json
    )
)

if not exist app.png (
    if exist app.png.in (
        echo Creating app.png from app.png.in
        copy /y app.png.in app.png
    )
)
rmdir /s /q %~dp0dist 2>nul
call npm config set registry https://registry.npmmirror.com/
call npm install
call npm run build
set ELECTRON_FILE_PATH=%~dp0dist\win-unpacked\go-boost-electron.exe
if "%EXPLORER_ELECTRON%"=="" (
    if exist %ELECTRON_FILE_PATH% (
        explorer.exe /select,"%ELECTRON_FILE_PATH%"
    ) else (
        explorer.exe /select,"%ELECTRON_FILE_PATH%\..\"
    )
)
popd

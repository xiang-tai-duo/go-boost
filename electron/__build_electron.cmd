@echo off
pushd %~dp0
set "APP_PNG=%~dp0src\electron\app.png"
set "PACKAGE_JSON=%~dp0src\electron\package.json"
if exist "%APP_PNG%" (
    copy "%APP_PNG%" %~dp0go-boost\electron\app.png /y 2>nul
)
if exist "%PACKAGE_JSON%" (
    copy "%PACKAGE_JSON%" %~dp0go-boost\electron\package.json /y 2>nul
)
call %~dp0go-boost\electron\__build.cmd
popd

@echo off
setlocal enabledelayedexpansion

echo Building upx.exe...
set GOWORK=off
go build -ldflags "-s -w" -o upx.exe .

if %errorlevel% equ 0 (
    echo Build succeeded: upx.exe
) else (
    echo Build failed
)

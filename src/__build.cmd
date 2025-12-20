@echo off
chcp 65001 >nul

set "OUTPUT_NAME=APP_NAME"
cd /d "%~dp0"

if not exist "bin" (
    mkdir "bin"
)

set "MINGW32_ENV=C:\Program Files\mingw32\mingwvars.bat"
set "MINGW64_ENV=C:\Program Files\mingw64\mingwvars.bat"

:: Build 32-bit version
set "GOOS=windows"
set "GOARCH=386"
set "CGO_ENABLED=1"
set "CC=gcc"
set "CXX=g++"

if exist "%MINGW32_ENV" (
    call "%MINGW32_ENV"
)

set "OUTPUT_32_EXE=bin\%OUTPUT_NAME%32.exe"
set "OUTPUT_32_DLL=bin\%OUTPUT_NAME%32.dll"

go build -o "%OUTPUT_32_EXE%"
go build -buildmode=c-shared -o "%OUTPUT_32_DLL%"

:: Build 64-bit version
set "GOOS=windows"
set "GOARCH=amd64"
set "CGO_ENABLED=1"
set "CC=gcc"
set "CXX=g++"

if exist "%MINGW64_ENV" (
    call "%MINGW64_ENV"
)

set "OUTPUT_64_EXE=bin\%OUTPUT_NAME%64.exe"
set "OUTPUT_64_DLL=bin\%OUTPUT_NAME%64.dll"

go build -o "%OUTPUT_64_EXE%"
go build -buildmode=c-shared -o "%OUTPUT_64_DLL%"

echo Build completed:
echo 32-bit EXE: %OUTPUT_32_EXE%
echo 32-bit DLL: %OUTPUT_32_DLL%
echo 64-bit EXE: %OUTPUT_64_EXE%
echo 64-bit DLL: %OUTPUT_64_DLL%

pause
exit /b 0
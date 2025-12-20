@echo off
set "APP_NAME=APP"
set "BUILD_TAGS=windows"
set "ELECTRON_DIST_EMBED=0"

setlocal enabledelayedexpansion
chcp 65001 >nul
if "%ELECTRON_DIST_EMBED%"=="0" (
    set "EXPLORER_ELECTRON=0"
    if exist "%~dp0__build_electron.cmd" (
        call "%~dp0__build_electron.cmd"
    )
)
pushd "%~dp0src"
set "BUILD_EXE=1"
set "BUILD_DLL=0"
set "BUILD_32BIT=0"
set "BUILD_64BIT=1"
set "CC=gcc"
set "CGO_ENABLED=1"
set "CXX=g++"
set "DIRBIN=bin\"
set "DIRPATH32=%DIRBIN%i386\"
set "DIRPATH64=%DIRBIN%amd64\"
set "DLLNAME=%APP_NAME%.dll"
set "DLLPATH32=%DIRPATH32%%DLLNAME%"
set "DLLPATH64=%DIRPATH64%%DLLNAME%"
set "EXENAME=%APP_NAME%.exe"
set "EXEPATH32=%DIRPATH32%%EXENAME%"
set "EXEPATH64=%DIRPATH64%%EXENAME%"
set "GOOS=windows"
set "MINGW32_ENV=C:\PROGRA~1\mingw32\mingwvars.bat"
set "MINGW64_ENV=C:\PROGRA~1\mingw64\mingwvars.bat"
rmdir /s /q "bin"
mkdir "bin"
if not exist "winres" (
    mkdir "winres"
)
if not exist "winres\icon.png" (
    copy "..\go-boost\winres\icon.png" "winres\icon.png"
)
if not exist "winres\winres.json" (
    copy "..\go-boost\winres\winres.json" "winres\winres.json"
)
go install github.com/tc-hib/go-winres@latest
go-winres make 2>nul
go mod tidy
if "%CGO_ENABLED%"=="1" (
    set LDFLAGS_EXE=-ldflags="-extldflags=-static -H=windowsgui"
    set LDFLAGS_DLL=-ldflags=-extldflags=-static
) else (
    set LDFLAGS_EXE=-ldflags="-H=windowsgui"
    set LDFLAGS_DLL=
)
if "%BUILD_TAGS%"=="" (
    set TAGS_PARAM=
) else (
    set TAGS_PARAM=-tags "%BUILD_TAGS%"
)
if "%BUILD_32BIT%"=="1" (
    set "GOARCH=386"
    if exist "%MINGW32_ENV%" (
        call "%MINGW32_ENV%" 1>nul
    )
    if "%BUILD_EXE%"=="1" (
        go build %LDFLAGS_EXE% %TAGS_PARAM% -o "%EXEPATH32%"
        if exist "app.exe.manifest" (
            copy "app.exe.manifest" "%DIRPATH32%%EXENAME%.manifest"
        ) else (
            copy "..\go-boost\winres\app.exe.manifest" "%DIRPATH32%%EXENAME%.manifest"
        )
        if "%ELECTRON_DIST_EMBED%"=="0" (
            xcopy "..\go-boost\electron\dist\win-unpacked\" "%DIRPATH32%dist\win-unpacked" /e /c /h /i /y
        )
    )
    if "%BUILD_DLL%"=="1" (
        go build %LDFLAGS_DLL% %TAGS_PARAM% -buildmode=c-shared -o "%DLLPATH32%"
    )
)
if "%BUILD_64BIT%"=="1" (
    set "GOARCH=amd64"
    if exist "%MINGW64_ENV%" (
        call "%MINGW64_ENV%" 1>nul
    )
    if "%BUILD_EXE%"=="1" (
        go build %LDFLAGS_EXE% %TAGS_PARAM% -o "%EXEPATH64%"
        if exist "app.exe.manifest" (
            copy "app.exe.manifest" "%DIRPATH64%%EXENAME%.manifest"
        ) else (
            copy "..\go-boost\winres\app.exe.manifest" "%DIRPATH64%%EXENAME%.manifest"
        )
        if "%ELECTRON_DIST_EMBED%"=="0" (
            xcopy "..\go-boost\electron\dist\win-unpacked\" "%DIRPATH64%dist\win-unpacked" /e /c /h /i /y
        )
    )
    if "%BUILD_DLL%"=="1" (
        go build %LDFLAGS_DLL% %TAGS_PARAM% -buildmode=c-shared -o "%DLLPATH64%"
    )
)
if "%BUILD_32BIT%"=="1" (
    if "%BUILD_EXE%"=="1" (
        echo 32-bit EXE: %EXEPATH32%
    )
    if "%BUILD_DLL%"=="1" (
        echo 32-bit DLL: %DLLPATH32%
    )
)
if "%BUILD_64BIT%"=="1" (
    if "%BUILD_EXE%"=="1" (
        echo 64-bit EXE: %EXEPATH64%
    )
    if "%BUILD_DLL%"=="1" (
        echo 64-bit DLL: %DLLPATH64%
    )
)
popd
explorer.exe %~dp0src\bin
exit /b 0

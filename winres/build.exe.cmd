@echo off
@pushd %~dp0

set OUTPUT_NAME=GenericApp
echo =========================================
echo Building %OUTPUT_NAME% EXE
echo =========================================

if not exist "%~dp0bin" (
    echo Creating bin directory...
    mkdir "%~dp0bin"
)

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

go install github.com/tc-hib/go-winres@latest
go-winres make 2>nul

del bin\%OUTPUT_NAME%.exe 2>nul

pushd "%~dp0.."
go build -o "%~dp0bin\%OUTPUT_NAME%.exe"
popd

echo =========================================
echo Build completed successfully!
echo =========================================
popd
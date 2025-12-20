@echo off
@pushd %~dp0

set OUTPUT_NAME=GenericLib
echo =========================================
echo Building %OUTPUT_NAME% DLL
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

del bin\%OUTPUT_NAME%.h 2>nul
del bin\%OUTPUT_NAME%.dll 2>nul

pushd "%~dp0.."
go build -buildmode=c-shared -o "%~dp0bin\%OUTPUT_NAME%.dll"
popd

copy "%~dp0bin\%OUTPUT_NAME%.dll" "%~dp0..\csharp\bin\x64\Release\%OUTPUT_NAME%.dll" /y 2>nul
copy "%~dp0bin\%OUTPUT_NAME%.h" "%~dp0..\csharp\bin\x64\Release\%OUTPUT_NAME%.h" /y 2>nul

echo =========================================
echo Build completed successfully!
echo =========================================
popd
@echo off

rem Create output directory
mkdir .tmp 2>nul

rem Iterate through all .go files
for %%f in (*.go) do (
    rem Get filename without extension
    set "filename=%%~nf"
    
    echo Compiling %%f...
    go build -o .tmp\%%~nf.exe "%%f"
    
    rem Check if compilation succeeded
    if %errorlevel% equ 0 (
        echo ✓ Successfully compiled %%~nf
    ) else (
        echo ✗ Failed to compile %%~nf
    )
)

echo.
echo Build completed! Executables are in .tmp\ directory.

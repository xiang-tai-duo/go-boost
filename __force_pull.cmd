@echo off
cd /d %~dp0
set "REMOTE="
for /f "tokens=*" %%i in ('git remote') do (
    if not defined REMOTE set "REMOTE=%%i"
)
if not defined REMOTE (
    echo Error: No git remotes found!
    exit /b 1
)
git fetch %REMOTE%
git checkout master
git reset --hard %REMOTE%/master
git clean -fd
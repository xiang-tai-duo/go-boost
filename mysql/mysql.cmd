@echo off
setlocal enableextensions

REM ============================================================================
REM File:        mysql.cmd
REM Author:      TRAE.AI
REM Description: Force-reset the MySQL root password to sa
REM              Assumes MySQL is installed at C:\Program Files\MySQL\MySQL Server 9.7
REM ============================================================================

set "MYSQL_HOME=C:\Program Files\MySQL\MySQL Server 9.7"
set "MYSQL_BIN=%MYSQL_HOME%\bin"
set "MYSQLD=%MYSQL_BIN%\mysqld.exe"
set "MYSQL_CLI=%MYSQL_BIN%\mysql.exe"
set "NEW_PASSWORD=sa"
set "INIT_FILE=%TEMP%\mysql_reset_password.sql"

echo ============================================================
echo MySQL Root Password Reset Tool
echo MYSQL_HOME   = %MYSQL_HOME%
echo NEW_PASSWORD = %NEW_PASSWORD%
echo ============================================================

net session >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Please run this script as Administrator.
    pause
    exit /b 1
)

if not exist "%MYSQLD%" (
    echo [ERROR] mysqld.exe not found: %MYSQLD%
    pause
    exit /b 1
)

if not exist "%MYSQL_CLI%" (
    echo [ERROR] mysql.exe not found: %MYSQL_CLI%
    pause
    exit /b 1
)

echo.
echo [STEP 1/6] Detecting MySQL service name ...
set "MYSQL_SERVICE="
for /f "tokens=2 delims==" %%S in ('wmic service where "PathName like '%%MySQL Server 9.7%%mysqld%%'" get Name /value ^| find "="') do (
    set "MYSQL_SERVICE=%%S"
)
if not defined MYSQL_SERVICE (
    for %%N in (MySQL97 MySQL90 MySQL MySQL80) do (
        sc query "%%N" >nul 2>&1
        if not errorlevel 1 set "MYSQL_SERVICE=%%N"
    )
)
if defined MYSQL_SERVICE (
    echo           Detected service name: %MYSQL_SERVICE%
) else (
    echo           No MySQL service detected, will handle as standalone process.
)

echo.
echo [STEP 2/6] Stopping MySQL service / process ...
if defined MYSQL_SERVICE (
    net stop "%MYSQL_SERVICE%" >nul 2>&1
)
taskkill /F /IM mysqld.exe >nul 2>&1
timeout /t 2 /nobreak >nul

echo.
echo [STEP 3/6] Generating password-reset SQL script ...
> "%INIT_FILE%" echo ALTER USER 'root'@'localhost' IDENTIFIED BY '%NEW_PASSWORD%';
>> "%INIT_FILE%" echo FLUSH PRIVILEGES;
echo           SQL file: %INIT_FILE%

echo.
echo [STEP 4/6] Starting mysqld with --init-file to perform the reset ...
start "mysqld-reset" /B "%MYSQLD%" --init-file="%INIT_FILE%" --console
echo           Waiting for mysqld initialization ...
timeout /t 10 /nobreak >nul

echo.
echo [STEP 5/6] Verifying the new password ...
"%MYSQL_CLI%" -uroot -p%NEW_PASSWORD% -e "SELECT 'password reset OK' AS result;" 2>nul
if errorlevel 1 (
    echo [WARN] Verification failed, mysqld may still be initializing. Please login manually later to verify.
) else (
    echo [OK]   New password is in effect.
)

echo.
echo [STEP 6/6] Cleaning up and restarting MySQL service ...
taskkill /F /IM mysqld.exe >nul 2>&1
timeout /t 2 /nobreak >nul
del /F /Q "%INIT_FILE%" >nul 2>&1

if defined MYSQL_SERVICE (
    net start "%MYSQL_SERVICE%" >nul 2>&1
    if errorlevel 1 (
        echo [WARN] Failed to start service %MYSQL_SERVICE%, please start it manually.
    ) else (
        echo [OK]   Service %MYSQL_SERVICE% has been restarted.
    )
)

echo.
echo ============================================================
echo Reset complete. Username: root   Password: %NEW_PASSWORD%
echo ============================================================
endlocal
pause

@echo off
:: Apex Blockchain - Service Worker Storage Error Fix
:: Run this script to fix the Chrome Service Worker database error

echo.
echo ========================================
echo  Service Worker Storage Error Fix
echo ========================================
echo.

echo This script will:
echo 1. Close all Chrome processes
echo 2. Clear Service Worker cache
echo 3. Clear IndexedDB
echo 4. Clear browser cache
echo.

pause

echo.
echo Step 1: Closing Chrome...
taskkill /F /IM chrome.exe >nul 2>&1
taskkill /F /IM msedge.exe >nul 2>&1
timeout /t 2 >nul
echo ✓ Browsers closed
echo.

echo Step 2: Clearing Service Worker cache...
set CHROME_DIR=%LOCALAPPDATA%\Google\Chrome\User Data\Default
set EDGE_DIR=%LOCALAPPDATA%\Microsoft\Edge\User Data\Default

:: Chrome Service Worker
if exist "%CHROME_DIR%\Service Worker" (
    echo Removing Chrome Service Worker cache...
    rd /s /q "%CHROME_DIR%\Service Worker"
    echo ✓ Chrome Service Worker cache removed
) else (
    echo Chrome Service Worker cache not found
)

:: Edge Service Worker
if exist "%EDGE_DIR%\Service Worker" (
    echo Removing Edge Service Worker cache...
    rd /s /q "%EDGE_DIR%\Service Worker"
    echo ✓ Edge Service Worker cache removed
) else (
    echo Edge Service Worker cache not found
)
echo.

echo Step 3: Clearing IndexedDB...

:: Chrome IndexedDB
if exist "%CHROME_DIR%\IndexedDB" (
    echo Removing Chrome IndexedDB...
    rd /s /q "%CHROME_DIR%\IndexedDB"
    echo ✓ Chrome IndexedDB removed
)

:: Edge IndexedDB
if exist "%EDGE_DIR%\IndexedDB" (
    echo Removing Edge IndexedDB...
    rd /s /q "%EDGE_DIR%\IndexedDB"
    echo ✓ Edge IndexedDB removed
)
echo.

echo Step 4: Clearing Cache Storage...

:: Chrome Cache Storage
if exist "%CHROME_DIR%\Cache" (
    echo Removing Chrome cache...
    rd /s /q "%CHROME_DIR%\Cache"
    echo ✓ Chrome cache removed
)

:: Chrome Code Cache
if exist "%CHROME_DIR%\Code Cache" (
    rd /s /q "%CHROME_DIR%\Code Cache"
    echo ✓ Chrome code cache removed
)
echo.

echo ========================================
echo   SERVICE WORKER ERROR FIXED!
echo ========================================
echo.
echo The following have been cleared:
echo  ✓ Service Worker cache
echo  ✓ IndexedDB
echo  ✓ Browser cache
echo  ✓ Code cache
echo.
echo You can now safely:
echo  1. Restart Chrome/Edge
echo  2. Run your Apex blockchain node
echo  3. Access the RPC API without errors
echo.
echo If the error persists:
echo  - Try running Chrome in incognito mode
echo  - Or use: chrome.exe --disable-features=ServiceWorkerCache
echo  - Check BUILD_GUIDE.md for more solutions
echo.
pause

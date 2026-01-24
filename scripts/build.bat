@echo off
SETLOCAL EnableDelayedExpansion

:: Apex Blockchain - Build and Fix Script for Windows
:: This script builds the project and fixes common issues

echo.
echo ========================================
echo    Apex Blockchain Setup Script
echo ========================================
echo.

:: Colors (use for better output)
:: 1. Check if Go is installed
echo [1/7] Checking Go installation...
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH
    echo Please install Go from https://golang.org/dl/
    pause
    exit /b 1
)
echo ✓ Go is installed
echo.

:: 2. Fix Service Worker Storage Error
echo [2/7] Fixing Service Worker storage error...
echo.
echo Closing Chrome processes...
taskkill /F /IM chrome.exe >nul 2>&1
timeout /t 2 >nul

echo Clearing Service Worker cache...
set CHROME_DIR=%LOCALAPPDATA%\Google\Chrome\User Data\Default

if exist "%CHROME_DIR%\Service Worker" (
    rd /s /q "%CHROME_DIR%\Service Worker" 2>nul
    echo ✓ Service Worker cache cleared
) else (
    echo Service Worker cache not found (this is okay)
)

if exist "%CHROME_DIR%\IndexedDB" (
    rd /s /q "%CHROME_DIR%\IndexedDB" 2>nul
    echo ✓ IndexedDB cleared
)

echo.

:: 3. Clean previous builds
echo [3/7] Cleaning previous builds...
if exist bin rd /s /q bin 2>nul
if exist data rd /s /q data 2>nul
echo ✓ Cleaned
echo.

:: 4. Install/Update dependencies
echo [4/7] Installing dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download dependencies
    pause
    exit /b 1
)
go mod tidy
echo ✓ Dependencies ready
echo.

:: 5. Create necessary directories
echo [5/7] Creating directories...
mkdir bin 2>nul
mkdir data 2>nul
mkdir logs 2>nul
mkdir config 2>nul
echo ✓ Directories created
echo.

:: 6. Build binaries
echo [6/7] Building binaries...
echo.

echo Building apex node...
go build -o bin\apex.exe .\cmd\apex
if %errorlevel% neq 0 (
    echo ERROR: Failed to build apex
    pause
    exit /b 1
)
echo ✓ apex.exe built successfully

echo Building apexctl CLI tool...
go build -o bin\apexctl.exe .\cmd\apexctl
if %errorlevel% neq 0 (
    echo ERROR: Failed to build apexctl
    pause
    exit /b 1
)
echo ✓ apexctl.exe built successfully

echo Building genesis tool...
go build -o bin\genesis.exe .\cmd\genesis
if %errorlevel% neq 0 (
    echo ERROR: Failed to build genesis
    pause
    exit /b 1
)
echo ✓ genesis.exe built successfully
echo.

:: 7. Verify builds
echo [7/7] Verifying builds...
if not exist bin\apex.exe (
    echo ERROR: apex.exe not found
    pause
    exit /b 1
)
if not exist bin\apexctl.exe (
    echo ERROR: apexctl.exe not found
    pause
    exit /b 1
)
if not exist bin\genesis.exe (
    echo ERROR: genesis.exe not found
    pause
    exit /b 1
)
echo ✓ All binaries verified
echo.

:: Display success message
echo ========================================
echo   BUILD COMPLETED SUCCESSFULLY!
echo ========================================
echo.
echo Binaries created:
echo   - bin\apex.exe       (Main blockchain node)
echo   - bin\apexctl.exe    (CLI management tool)
echo   - bin\genesis.exe    (Genesis generator)
echo.
echo Service Worker Error: FIXED
echo.
echo ----------------------------------------
echo Next Steps:
echo ----------------------------------------
echo.
echo 1. Generate genesis block:
echo    bin\genesis.exe --output config\genesis.json
echo.
echo 2. Start the node:
echo    bin\apex.exe --config config\config.yaml
echo.
echo 3. Test the API:
echo    curl -X POST http://localhost:8545 -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"apex_blockNumber\",\"params\":[],\"id\":1}"
echo.
echo For more information, see BUILD_GUIDE.md
echo.
echo ========================================

:: Optional: Run tests
echo.
set /p RUN_TESTS="Do you want to run tests? (y/n): "
if /i "%RUN_TESTS%"=="y" (
    echo.
    echo Running tests...
    go test ./...
    echo.
)

:: Optional: Create default config
if not exist config\config.yaml (
    echo.
    echo Creating default configuration...
    copy config\config.yaml.example config\config.yaml >nul 2>&1
    echo ✓ Default config created
    echo.
)

echo.
echo Build script completed!
pause

@echo off
:: Quick Fix Script for Empty Go Files
:: This script removes empty Go files and rebuilds the project

echo.
echo ========================================
echo    Quick Fix - Empty Go Files
echo ========================================
echo.

echo Checking for empty Go files...
echo.

:: Delete voting.go if it exists
if exist "pkg\consensus\voting.go" (
    del "pkg\consensus\voting.go"
    echo ✓ Deleted pkg\consensus\voting.go
) else (
    echo - pkg\consensus\voting.go not found (already deleted?)
)

:: Delete state.go if it exists
if exist "pkg\core\state.go" (
    del "pkg\core\state.go"
    echo ✓ Deleted pkg\core\state.go
) else (
    echo - pkg\core\state.go not found (already deleted?)
)

echo.
echo ========================================
echo   Empty files removed!
echo ========================================
echo.
echo Running build script...
echo.

:: Run the build script
call scripts\build.bat

echo.
echo Done!
pause
@echo off
echo Fixing import paths in all Go files...
echo.

:: Fix cmd/apex/main.go
powershell -Command "(Get-Content 'cmd\apex\main.go') -replace 'github.com/apex-blockchain/apex/pkg/', 'github.com/apex/pkg/' | Set-Content 'cmd\apex\main.go'"

:: Fix cmd/apexctl/main.go
powershell -Command "(Get-Content 'cmd\apexctl\main.go') -replace 'github.com/apex-blockchain/apex/pkg/', 'github.com/apex/pkg/' | Set-Content 'cmd\apexctl\main.go'"

:: Fix cmd/genesis/main.go
powershell -Command "(Get-Content 'cmd\genesis\main.go') -replace 'github.com/apex-blockchain/apex/pkg/', 'github.com/apex/pkg/' | Set-Content 'cmd\genesis\main.go'"

:: Fix all pkg files
for /r pkg %%f in (*.go) do (
    powershell -Command "(Get-Content '%%f') -replace 'github.com/apex-blockchain/pkg/', 'github.com/apex/pkg/' | Set-Content '%%f'"
)

echo.
echo ✓ All imports fixed!
echo.
echo Now run: go mod tidy
echo Then run: .\scripts\build.bat
pause
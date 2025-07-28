@echo off
echo ===== SENTINEL Release Builder =====

REM Check if Go is installed
where go >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo Error: Go not found. Make sure Go is installed and in your PATH.
    exit /b 1
)

REM Check if GoReleaser is installed
where goreleaser >nul 2>nul
if %ERRORLEVEL% neq 0 (
    echo GoReleaser not found. Installing GoReleaser...
    go install github.com/goreleaser/goreleaser@latest
    if %ERRORLEVEL% neq 0 (
        echo Error: Failed to install GoReleaser.
        exit /b 1
    )
    echo GoReleaser installed successfully.
)

REM Clean dist folder if it exists
if exist dist (
    echo Cleaning dist folder...
    rmdir /s /q dist
)

REM Run GoReleaser in snapshot mode (without uploading to GitHub)
echo Creating local release...
goreleaser release --snapshot --clean

if %ERRORLEVEL% neq 0 (
    echo Error: Failed to create release.
    exit /b 1
)

echo.
echo ===== Release created successfully! =====
echo Binaries available in dist\ folder
echo.

REM Display list of generated files
echo Generated files:
dir /b dist\

pause
@echo off
REM LogTrawl Build Script
REM Supports building for Windows, Linux, and macOS platforms

setlocal enabledelayedexpansion

echo ========================================
echo LogTrawl Application Build Script
echo ========================================

REM Check if Wails is installed
where wails >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Wails command not found. Please install Wails v2:
    echo Run: go install github.com/wailsapp/wails/v2/cmd/wails@latest
    exit /b 1
)

REM Show help information
if "%1"=="" (
    echo.
    echo Usage:
    echo   build.bat [platform]
    echo.
    echo Platform options:
    echo   win     - Build Windows application
    echo   linux   - Build Linux application
    echo   mac     - Build macOS application
    echo   all     - Build applications for all platforms
    echo   clean   - Clean build directory
    echo.
    echo Examples:
    echo   build.bat win
    echo   build.bat all
    echo.
    exit /b 0
)

REM Get current directory
set "PROJECT_DIR=%~dp0"
set "BUILD_DIR=%PROJECT_DIR%dist"

REM Create build directory
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

REM Handle command arguments
if "%1"=="clean" (
    echo Cleaning build directory...
    if exist "%BUILD_DIR%" (
        rmdir /s /q "%BUILD_DIR%"
        echo Build directory cleaned
    ) else (
        echo Build directory does not exist
    )
    exit /b 0
)

if "%1"=="win" (
    call :build_windows
    exit /b %errorlevel%
)

if "%1"=="linux" (
    call :build_linux
    exit /b %errorlevel%
)

if "%1"=="mac" (
    call :build_macos
    exit /b %errorlevel%
)

if "%1"=="all" (
    call :build_windows
    if !errorlevel! neq 0 exit /b !errorlevel!
    
    echo.
    echo ========================================
    echo Windows platform built successfully!
    echo ========================================
    echo WARNING: Cross-compilation to Linux and macOS is not supported on Windows.
    echo To build for Linux or macOS, please run this script on the respective platforms.
    echo ========================================
    exit /b 0
)

echo Unknown platform option: %1
echo Use "build.bat" to see help information
exit /b 1

REM Build Windows application
:build_windows
echo.
echo ========================================
echo Building Windows application
echo ========================================
cd /d "%PROJECT_DIR%"

REM Build Windows application
wails build -platform windows/amd64
if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"
if exist "build\bin\LogTrawl.exe" (
    move "build\bin\LogTrawl.exe" "%BUILD_DIR%\LogTrawl-windows-amd64.exe"
)

if %errorlevel% equ 0 (
    echo Windows application built successfully: %BUILD_DIR%\LogTrawl-windows-amd64.exe
) else (
    echo Windows application build failed
    exit /b %errorlevel%
)
exit /b 0

REM Build Linux application
:build_linux
echo.
echo ========================================
echo Building Linux application
echo ========================================
cd /d "%PROJECT_DIR%"

REM Check if running on Linux (cross-compilation to Linux is not supported on Windows)
echo WARNING: Cross-compilation to Linux is not supported on Windows.
echo To build for Linux, please run this script on a Linux machine.
echo.
echo For more information, visit: https://github.com/sponsors/leaanthony
exit /b 1

REM Build macOS application
:build_macos
echo.
echo ========================================
echo Building macOS application
echo ========================================
cd /d "%PROJECT_DIR%"

REM Check if running on macOS (cross-compilation to macOS is not supported on Windows)
echo WARNING: Cross-compilation to macOS is not supported on Windows.
echo To build for macOS, please run this script on a macOS machine.
echo.
echo For more information, visit: https://github.com/sponsors/leaanthony
exit /b 1
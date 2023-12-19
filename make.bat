@echo off

setlocal

set dist=dist
set target=wzry-skin
set build=go.exe build -ldflags="-s -w"

if "%~1" == "clean" goto lbl-clean
if "%~1" == "windows" goto lbl-windows
if "%~1" == "linux" goto lbl-linux
if "%~1" == "darwin" goto lbl-darwin
goto lbl-build

:lbl-clean
rd /s /q %dist%
exit /b

:lbl-build
call :lbl-windows
call :lbl-linux
call :lbl-darwin
exit /b

:lbl-cross-platform
rem cross platform compile
if not defined CGO_ENABLED set CGO_ENABLED=0
if not defined GOARCH set GOARCH=amd64
exit /b

:lbl-windows
rem build for Windows
%build% -o %dist%/%target%-windows.exe
exit /b

:lbl-linux
rem build for Linux
call :lbl-cross-platform
set GOOS=linux
%build% -o %dist%/%target%-linux
exit /b

:lbl-darwin
rem build for macOS
call :lbl-cross-platform
set GOOS=darwin
%build% -o %dist%/%target%-darwin
exit /b

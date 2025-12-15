@echo off

pushd %~dp0

setlocal

set dist=dist
if not exist %dist% mkdir %dist%

set gobuild=go build -ldflags="-s -w" -trimpath -o

set output=wzry-skin

%gobuild% %dist%\%output%.exe

if "%~1" neq "all" exit /b

set GOARCH=amd64
set GOOS=linux
%gobuild% %dist%\%output%-linux

set GOARCH=arm64
set GOOS=darwin
%gobuild% %dist%\%output%-macos

popd

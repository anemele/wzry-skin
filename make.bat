@echo off

pushd %~dp0

setlocal

set output=wzry-skin
set dist=dist
if not exist %dist% mkdir %dist%

go build -ldflags="-s -w" -trimpath -o %dist%\%output%.exe

popd

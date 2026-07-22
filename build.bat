@echo off
title MapRenderer Installer


echo Building MapRenderer V2...

where go >nul 2>nul
if errorlevel 1 (
    echo.
    echo Go is not installed or not in PATH.
    pause
    exit /b 1
)

set CGO_ENABLED=1

go build -o Python_GUI\bin\maprenderer.exe ./cmd/maprenderer

if errorlevel 1 (
    echo.
    echo Build failed.
    echo Maybe try manual approach?
    echo set CGO_ENABLED=1
    echo Try ... go build -o Python_GUI/bin/maprenderer.exe ./cmd/maprenderer
    pause
    exit /b 1
)

echo.
echo Build completed successfully.
echo Output:
echo Python_GUI\bin\maprenderer.exe
echo Note: If you see some nodes missing kindly check probe.go.
pause
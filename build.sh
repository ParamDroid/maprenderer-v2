#!/bin/bash


echo "Building MapRenderer V2..."


if ! command -v go >/dev/null 2>&1; then
    echo "Go is not installed."
    exit 1
fi

export CGO_ENABLED=1

go build -o Python_GUI/bin/maprenderer ./cmd/maprenderer

if [ $? -ne 0 ]; then
    echo "Build failed."
    echo "Maybe try manual approach?"
    echo "export CGO_ENABLED=1"
    echo "Try ... go build -o Python_GUI/bin/maprenderer ./cmd/maprenderer "
    exit 1
fi

chmod +x Python_GUI/bin/maprenderer

echo
echo "Build complete."
echo "Output:"
echo "Python_GUI/bin/maprenderer"
echo "Note: If you see some nodes missing,check probe.go."
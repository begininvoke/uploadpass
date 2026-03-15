#!/bin/bash

APP="uploadpass"
OUT="build"

rm -rf "$OUT"
mkdir -p "$OUT"

success=0
fail=0

build() {
    local goos=$1 goarch=$2 cc=$3
    local name="${APP}-${goos}-${goarch}"
    [ "$goos" = "windows" ] && name="${name}.exe"

    printf "  %-40s" "$name"
    CGO_ENABLED=1 GOOS="$goos" GOARCH="$goarch" CC="$cc" \
        go build -ldflags="-s -w" -o "${OUT}/${name}" . 2>/dev/null

    if [ $? -eq 0 ]; then
        size=$(du -h "${OUT}/${name}" | cut -f1 | xargs)
        echo "OK  ($size)"
        ((success++))
    else
        echo "SKIP"
        ((fail++))
    fi
}

echo "========================================="
echo " Building ${APP}"
echo "========================================="
echo ""

echo "[macOS]"
build darwin  amd64 "gcc"
build darwin  arm64 "gcc"

echo ""
echo "[Linux]"
build linux amd64 "x86_64-linux-musl-gcc"
build linux arm64 "aarch64-linux-musl-gcc"
build linux arm   "arm-linux-musleabihf-gcc"

echo ""
echo "[Windows]"
build windows amd64 "x86_64-w64-mingw32-gcc"
build windows 386   "i686-w64-mingw32-gcc"

echo ""
echo "========================================="
echo " Done: ${success} built, ${fail} skipped"
echo "========================================="
echo ""

if [ $success -gt 0 ]; then
    echo "Output:"
    ls -lh "$OUT"/
fi

if [ $fail -gt 0 ]; then
    echo ""
    echo "Skipped builds need cross-compilers. Install with:"
    echo "  brew install FiloSottile/musl-cross/musl-cross    # Linux"
    echo "  brew install mingw-w64                            # Windows"
fi

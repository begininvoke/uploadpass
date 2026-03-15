#!/bin/bash

APP="uploadpass"
OUT="build"

rm -rf "$OUT"
mkdir -p "$OUT"

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm64"
    "linux/386"
    "linux/arm"
    "windows/amd64"
    "windows/arm64"
    "windows/386"
)

success=0

echo "========================================="
echo " Building ${APP}"
echo "========================================="
echo ""

for platform in "${platforms[@]}"; do
    GOOS="${platform%/*}"
    GOARCH="${platform#*/}"
    name="${APP}-${GOOS}-${GOARCH}"
    [ "$GOOS" = "windows" ] && name="${name}.exe"

    printf "  %-40s" "$name"
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
        go build -ldflags="-s -w" -o "${OUT}/${name}" .

    if [ $? -eq 0 ]; then
        size=$(du -h "${OUT}/${name}" | cut -f1 | xargs)
        echo "OK  ($size)"
        ((success++))
    else
        echo "FAIL"
    fi
done

echo ""
echo "========================================="
echo " Done: ${success}/${#platforms[@]} built"
echo "========================================="
echo ""
ls -lh "$OUT"/

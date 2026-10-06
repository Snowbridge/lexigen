#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

DIST="${DIST:-dist}"
LDFLAGS="-s -w"

usage() {
	cat <<'EOF'
Usage: scripts/build.sh <target>

Targets:
  all           Windows + Linux amd64
  local         binary for this OS/arch
  windows       dist/lexigen-windows-amd64.exe
  linux         dist/lexigen-linux-amd64
  linux-arm64   dist/lexigen-linux-arm64
  test          go test ./...
  clean         remove dist/
  help          this message
EOF
}

mkdir -p "$DIST"

go_build() {
	local goos=$1 goarch=$2 out=$3
	GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 go build -ldflags="$LDFLAGS" -o "$out" .
}

target="${1:-help}"

case "$target" in
help|-h|--help)
	usage
	;;
all)
	go_build windows amd64 "$DIST/lexigen-windows-amd64.exe"
	go_build linux amd64 "$DIST/lexigen-linux-amd64"
	echo "built: $DIST/lexigen-windows-amd64.exe $DIST/lexigen-linux-amd64"
	;;
local)
	out="$DIST/lexigen"
	case "$(uname -s)" in
	MINGW* | MSYS* | CYGWIN*) out="$DIST/lexigen.exe" ;;
	esac
	go build -ldflags="$LDFLAGS" -o "$out" .
	echo "built: $out"
	;;
windows)
	go_build windows amd64 "$DIST/lexigen-windows-amd64.exe"
	echo "built: $DIST/lexigen-windows-amd64.exe"
	;;
linux)
	go_build linux amd64 "$DIST/lexigen-linux-amd64"
	echo "built: $DIST/lexigen-linux-amd64"
	;;
linux-arm64)
	go_build linux arm64 "$DIST/lexigen-linux-arm64"
	echo "built: $DIST/lexigen-linux-arm64"
	;;
test)
	go test ./...
	;;
clean)
	rm -rf "$DIST"
	echo "removed $DIST/"
	;;
*)
	echo "unknown target: $target" >&2
	usage >&2
	exit 1
	;;
esac

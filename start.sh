#!/usr/bin/env bash
# Installs what grabb needs, builds it and opens the dashboard.
# Usage: ./start.sh [flags]   e.g. ./start.sh -out ~/Videos
set -euo pipefail
cd "$(dirname "$0")"

LOCAL="$HOME/.local"
export PATH="$LOCAL/bin:$LOCAL/go/bin:$HOME/.deno/bin:$PATH"
mkdir -p "$LOCAL/bin"

has() { command -v "$1" >/dev/null 2>&1; }
say() { printf '%s\n' "$*"; }

# go_ok is true when Go 1.25 or newer is on PATH.
go_ok() {
  has go || return 1
  local v; v=$(go env GOVERSION 2>/dev/null | sed 's/^go//')
  [ "$(printf '%s\n1.25\n' "$v" | sort -V | head -1)" = "1.25" ]
}

install_go() {
  local os arch version
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  case "$(uname -m)" in x86_64) arch=amd64 ;; aarch64 | arm64) arch=arm64 ;; *) arch=$(uname -m) ;; esac
  version=$(curl -fsSL "https://go.dev/VERSION?m=text" | head -1)
  say "Installing $version..."
  rm -rf "$LOCAL/go"
  curl -fsSL "https://go.dev/dl/$version.$os-$arch.tar.gz" | tar -xz -C "$LOCAL"
}

if [ "$(uname -s)" = "Darwin" ]; then
  has brew || { say "grabb needs Homebrew on macOS: https://brew.sh"; exit 1; }
  missing=()
  go_ok || missing+=(go)
  for t in yt-dlp ffmpeg deno; do has "$t" || missing+=("$t"); done
  has aria2c || missing+=(aria2)
  [ ${#missing[@]} -eq 0 ] || { say "Installing ${missing[*]}..."; brew install "${missing[@]}"; }
else
  if ! has ffmpeg || ! has aria2c || ! has curl; then
    if has apt-get; then sudo apt-get update -qq && sudo apt-get install -y -qq ffmpeg aria2 curl unzip
    elif has dnf; then sudo dnf install -y -q ffmpeg aria2 curl unzip
    elif has pacman; then sudo pacman -S --noconfirm --needed ffmpeg aria2 curl unzip
    else say "Install ffmpeg and aria2 with your package manager, then run this again."; exit 1
    fi
  fi
  go_ok || install_go
  if ! has yt-dlp; then
    say "Installing yt-dlp..."
    curl -fsSL https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o "$LOCAL/bin/yt-dlp"
    chmod +x "$LOCAL/bin/yt-dlp"
  fi
  has deno || { say "Installing deno..."; curl -fsSL https://deno.land/install.sh | sh -s -- -y >/dev/null; }
fi

say "Building grabb..."
go build -o bin/grabb ./cmd/grabb
exec ./bin/grabb "$@"

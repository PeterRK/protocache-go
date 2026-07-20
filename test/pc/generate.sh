#!/usr/bin/env bash
set -euo pipefail

repo_root=$(CDPATH= cd -- "$(dirname -- "$0")/../.." && pwd)
cd "$repo_root"

tmp_dir=$(mktemp -d "${TMPDIR:-/tmp}/protocache-go-generate.XXXXXX")
trap 'rm -rf "$tmp_dir"' EXIT

go build -o "$tmp_dir/protoc-gen-pcgo" ./cmd/protoc-gen-pcgo
protoc \
  --plugin="protoc-gen-pcgo=$tmp_dir/protoc-gen-pcgo" \
  --pcgo_out=test/pc \
  '--pcgo_opt=extra,relative,Mtest/test.proto=github.com/peterrk/protocache-go/test/pc;pc' \
  test/test.proto

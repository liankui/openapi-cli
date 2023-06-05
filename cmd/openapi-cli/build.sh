#!/bin/bash
set -ex
pwd=$(
  cd "$(dirname "$0")"
  pwd
)
version=$(cat "$(dirname "$pwd")/.."/version)
GOOSs=(darwin linux windows)
GOARCHs=(amd64)

for os in "${GOOSs[@]}"; do
  for arch in "${GOARCHs[@]}"; do
    GOOS=$os GOARCH=$arch go build -ldflags "-w -s" -o openapi-cli-"$version"-"$os"-"$arch"
  done
done

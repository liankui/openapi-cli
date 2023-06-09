#!/bin/bash
version=$(cat version)
GOOSs=(darwin linux windows)
GOARCHs=(amd64)

set -ex
for os in "${GOOSs[@]}"; do
  for arch in "${GOARCHs[@]}"; do
    GOOS=$os GOARCH=$arch go build -ldflags "-w -s" -o openapi-cli-"$version"-"$os"-"$arch"
  done
done

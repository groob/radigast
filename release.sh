#!/bin/bash

VERSION="0.0.1"

echo "Building Radigast version $VERSION"

mkdir -p pkg

build() {
  echo -n "=> $1-$2: "
  GOOS=$1 GOARCH=$2 go build -o pkg/radigast-$1-$2 -ldflags "-X main.Version $VERSION" ./cmd/radigast/radigast.go
  du -h pkg/telegraf-$1-$2
}

build "darwin" "amd64"
build "linux" "amd64"

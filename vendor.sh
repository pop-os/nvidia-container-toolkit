#!/bin/bash

set -x

GOLANG=go1.19.1.linux-amd64.tar.gz

export GOCACHE="${PWD}/gocache"
export GOPATH="${PWD}/gopath"
export PATH="$PATH:${PWD}/go/bin"
# export GO111MODULE=off

curl -L https://go.dev/dl/${GOLANG} --output ${GOLANG}
tar -C ${PWD} -xzf ${GOLANG}

rm -rf $GOCACHE $GOPATH vendor
mkdir -p $GOCACHE $GOPATH
go mod vendor
# go get -d github.com/NVIDIA/nvidia-container-toolkit/cmd/nvidia-container-runtime
# go get -d github.com/NVIDIA/nvidia-container-toolkit/cmd/nvidia-container-runtime-hook
# go get -d github.com/NVIDIA/nvidia-container-toolkit/cmd/nvidia-ctk
find gopath -name .git -exec rm -rf {} +
tar pcf vendor.tar gocache gopath
go clean -modcache
rm -rf ${PWD}/go ${GOLANG}
#!/bin/bash

# Builds on launchpad can't access network

export GOCACHE="$PWD/gocache"
export GOPATH="$PWD/gopath"
# Needed for bionic build
export GO111MODULE=off

rm -rf gopath
mkdir -p gopath/src/github.com/NVIDIA/nvidia-container-toolkit
ln -s "../../../../../pkg" gopath/src/github.com/NVIDIA/nvidia-container-toolkit
go get -d github.com/NVIDIA/nvidia-container-toolkit/pkg
find gopath -name .git -exec rm -rf {} +

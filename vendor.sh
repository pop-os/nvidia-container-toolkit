#!/bin/bash

export GOCACHE="${PWD}/gocache"
export GOPATH="${PWD}/gopath"
export GO111MODULE=off

rm -rf $GOCACHE $GOPATH
mkdir -p $GOCACHE $GOPATH
go get -d github.com/NVIDIA/nvidia-container-toolkit/cmd/nvidia-container-toolkit
find gopath -name .git -exec rm -rf {} +

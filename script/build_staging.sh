#!/bin/bash

set -ex
BinName=memnexus

# 先编译
rm -rf $BinName
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BinName ./cmd/main.go
scp $BinName $SSH_STAGING_USER@$SSH_STAGING_HOST:~
rm $BinName
echo "build $BinName success"

# todo: 再写个运行文件上去?
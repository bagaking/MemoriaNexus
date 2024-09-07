#!/bin/bash

set -ex
BinName=memnexus

# 先编译
rm -rf $BinName
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $BinName ./cmd/main.go

# 通过 app.sh kill 来关闭进程
ssh $SSH_STAGING_USER@$SSH_STAGING_HOST "~/app.sh kill"

# 上传
scp $BinName $SSH_STAGING_USER@$SSH_STAGING_HOST:~
rm $BinName

echo "build $BinName success"

# 重启
ssh $SSH_STAGING_USER@$SSH_STAGING_HOST "~/app.sh start"



# todo: 再写个运行文件上去?
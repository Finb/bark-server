#!/usr/bin/env bash

set -e

GO111MODULE="on"
GOPROXY="https://goproxy.cn"
TARGET_DIR="dist"
PLATFORMS="darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/386 windows/amd64 freebsd/386 freebsd/amd64"

rm -rf ${TARGET_DIR}
mkdir ${TARGET_DIR}

for pl in ${PLATFORMS}; do 
    GOOS=$(echo ${pl} | cut -d'/' -f1)
    GOARCH=$(echo ${pl} | cut -d'/' -f2)
    echo "build => ${TARGET_DIR}/bark-server_${GOOS}_${GOARCH}"
    go build -trimpath -o ${TARGET_DIR}/bark-server_${GOOS}_${GOARCH} \
            -ldflags    "-X 'main.Version=${BUILD_VERSION}' \
                        -X 'main.BuildDate=${BUILD_DATE}' \
                        -X 'main.CommitID=${COMMIT_SHA1}'\
                        -w -s"
done

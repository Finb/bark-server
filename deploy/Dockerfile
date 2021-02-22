FROM golang:1.16-alpine3.12 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

COPY . /go/src/github.com/finb/bark-server

WORKDIR /go/src/github.com/finb/bark-server

RUN set -ex \
    && apk add git gcc libc-dev \
    && BUILD_VERSION=`cat version` \
    && BUILD_DATE=`date "+%F %T"` \
    && COMMIT_SHA1=`git rev-parse HEAD` \
    && go install -trimpath -ldflags \
            "-X 'main.version=${BUILD_VERSION}' \
             -X 'main.buildDate=${BUILD_DATE}' \
             -X 'main.commitID=${COMMIT_SHA1}' \
             -w -s"

FROM alpine:3.12

LABEL maintainer="mritd <mritd@linux.com>"

RUN set -ex \
    && apk upgrade \
    && apk add ca-certificates bash

COPY --from=builder /go/bin/bark-server /usr/local/bin/bark-server

VOLUME /data

EXPOSE 8080

# Allow bark-server to run in Prefork mode
# see also: https://docs.gofiber.io/api/fiber#config
SHELL ["/bin/bash", "-c"]

CMD ["bark-server"]

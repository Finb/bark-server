FROM golang:1.13.5-alpine3.11 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
ENV GOSUMDB sum.golang.google.cn

COPY .. /go/src/github.com/finb/bark-server

WORKDIR /go/src/github.com/finb/bark-server

RUN set -ex \
    && apk add git \
    && BUILD_VERSION=`cat version` \
    && BUILD_DATE=`date "+%F %T"` \
    && COMMIT_SHA1=`git rev-parse HEAD` \
    && go install -ldflags  "-X 'main.Version=${BUILD_VERSION}' \
                             -X 'main.BuildDate=${BUILD_DATE}' \
                             -X 'main.CommitID=${COMMIT_SHA1}'"

FROM alpine:3.11

LABEL maintainer="mritd <mritd@linux.com>"

RUN apk upgrade \
    && apk add ca-certificates

COPY --from=builder /go/bin/bark-server /usr/local/bin/bark-server

VOLUME /data

EXPOSE 8080

CMD ["bark-server"]

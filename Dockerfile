FROM golang:1.13.1-alpine3.10 AS builder

ENV GO111MODULE on
ENV GOPROXY https://goproxy.io
ENV GOSUMDB sum.golang.google.cn

COPY . /go/src/github.com/finb/bark-server

WORKDIR /go/src/github.com/finb/bark-server

RUN go install

FROM alpine:3.10

LABEL maintainer="mritd <mritd1234@gmail.com>"

RUN apk upgrade \
    && apk add ca-certificates

COPY --from=builder /go/bin/bark-server /usr/local/bin/bark-server

VOLUME /data

EXPOSE 8080

CMD ["bark-server"]

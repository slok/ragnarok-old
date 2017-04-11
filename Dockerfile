FROM golang:1.8-alpine


RUN apk --update add musl-dev gcc tar git bash wget && rm -rf /var/cache/apk/*

# Create user
ARG uid=1000
ARG gid=1000
RUN addgroup -g $gid ragnarok
RUN adduser -D -u $uid -G ragnarok ragnarok

RUN mkdir -p /go/src/github.com/slok/ragnarok/
RUN chown -R ragnarok:ragnarok /go

WORKDIR /go/src/github.com/slok/ragnarok/

USER ragnarok

FROM slok/ragnarok_base:latest
USER root
RUN set -eux && \
    apk add --no-cache \
        g++ \
        protobuf

USER ragnarok

# golang dependencies
RUN go get github.com/Masterminds/glide && \
    go get github.com/vektra/mockery/.../ && \
    go get github.com/golang/protobuf/protoc-gen-go



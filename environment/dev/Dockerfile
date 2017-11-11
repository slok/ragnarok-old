FROM slok/ragnarok_base:latest
USER root

ENV PROTOBUF_VERSION=3.1.0
ENV GOOGLEPROTO_PATH=/usr/local/lib/proto

RUN set -eux && \
    apk add --update \
        g++ \
        curl \
        protobuf==${PROTOBUF_VERSION}-r1

# Install google proto files
RUN curl -L https://github.com/google/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip -o /tmp/pb.zip && \
    unzip /tmp/pb.zip -d /tmp && \
    mkdir -p ${GOOGLEPROTO_PATH} && \
    mv /tmp/include/google ${GOOGLEPROTO_PATH} && \
    chmod -R 755 ${GOOGLEPROTO_PATH} && \
    rm -rf /tmp/*


USER ragnarok

# golang dependencies
RUN go get github.com/golang/dep/cmd/dep && \
    go get github.com/vektra/mockery/.../ && \
    #go get github.com/golang/protobuf/protoc-gen-go && \
    go get github.com/gogo/protobuf/protoc-gen-gofast



FROM slok/ragnarok_base:latest

USER root

# Prepare
WORKDIR /go/src/github.com/slok/ragnarok/
RUN mkdir -p /bin
COPY . ./

# Build
RUN ./build.sh /bin/ragnarok
RUN chmod 755 /bin/ragnarok

# Clean up
WORKDIR /
RUN rm -rf /go/src/*


EXPOSE 9222

ENTRYPOINT [ "/bin/ragnarok" ]
CMD        [ "--help"]

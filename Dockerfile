##
# BUILD CONTAINER
##

FROM goreleaser/goreleaser:v0.148.0 as builder

WORKDIR /build

COPY . .
RUN \
apk add --no-cache make ca-certificates ;\
make build-linux-amd64

##
# RELEASE CONTAINER
##

FROM busybox:1.32-glibc

WORKDIR /

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/dist/s5_linux_amd64/s5 /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

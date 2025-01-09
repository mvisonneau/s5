##
# BUILD CONTAINER
##

FROM alpine:3.21@sha256:56fa17d2a7e7f168a043a2712e63aed1f8543aeafdcee47c58dcffe38ed51099 as certs

RUN \
  apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.37.0-glibc@sha256:c598938e58d0efcc5a01efe9059d113f22970914e05e39ab2a597a10f9db9bdc

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY s5 /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

##
# BUILD CONTAINER
##

FROM alpine:3.17@sha256:a6063e988bcd597b4f1f7cfd4ec38402b02edd0c79250f00c9e14dc1e94bebbc as certs

RUN \
  apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.35.0-glibc@sha256:b4899072f500eabf2504e7d0348955b46cd3f60dcbb3bc97ff56e5ef793263f7

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY s5 /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

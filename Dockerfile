##
# BUILD CONTAINER
##

FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN \
  apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.37.0-glibc@sha256:210ce53959959e79523b8cb0f0bb1cf1c49bf9747cdedb47db1cf0db8e642f61

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY s5 /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

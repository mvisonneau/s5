##
# BUILD CONTAINER
##

FROM alpine:3.22@sha256:8a1f59ffb675680d47db6337b49d22281a139e9d709335b492be023728e11715 as certs

RUN \
  apk add --no-cache ca-certificates

##
# RELEASE CONTAINER
##

FROM busybox:1.37.0-glibc@sha256:bd606c263abed91a141187b92fdb54b87bbc39cfb9068f96ad84196a36963103

WORKDIR /

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY s5 /usr/local/bin/

# Run as nobody user
USER 65534

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

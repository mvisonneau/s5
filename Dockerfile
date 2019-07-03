##
# BUILD CONTAINER
##

FROM golang:1.12.5 as builder

WORKDIR /build

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make build-docker

##
# RELEASE CONTAINER
##

FROM busybox:1.31-glibc

WORKDIR /

COPY --from=builder /build/s5 /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

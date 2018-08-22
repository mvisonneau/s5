##
# BUILD CONTAINER
##

FROM golang:1.10 as builder

WORKDIR /go/src/github.com/mvisonneau/s5

COPY Makefile .
RUN \
make setup

COPY . .
RUN \
make deps ;\
make build-docker

##
# RELEASE CONTAINER
##

FROM busybox:1.29

WORKDIR /usr/local/bin

COPY --from=builder /go/src/github.com/mvisonneau/s5/s5 /usr/local/bin

ENTRYPOINT ["/usr/local/bin/s5"]
CMD [""]

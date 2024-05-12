ARG GO_VERSION=1.22
ARG COREDNS_VERSION=1.11.3

FROM golang:${GO_VERSION}-alpine as builder

ARG COREDNS_VERSION

RUN apk add --no-cache git make

WORKDIR /coredns

RUN git clone -b v${COREDNS_VERSION} --depth 1 https://github.com/coredns/coredns.git .

RUN sed -i '/cache:cache/i stats:github.com/GustavoKatel/coredns-stats' plugin.cfg

ADD . /stats

RUN echo "replace github.com/GustavoKatel/coredns-stats => /stats" >> go.mod && \
    go get github.com/GustavoKatel/coredns-stats

RUN go generate && make

FROM alpine:latest

COPY --from=builder /coredns/coredns /coredns/coredns

ENTRYPOINT ["/coredns/coredns"]

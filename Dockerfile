ARG GO_VERSION=1.23.2
ARG COREDNS_VERSION=1.11.3

FROM golang:${GO_VERSION}-alpine as builder

ARG COREDNS_VERSION

RUN apk add --no-cache git make

WORKDIR /coredns

RUN git clone --depth 1 --branch v${COREDNS_VERSION} https://github.com/coredns/coredns.git .

RUN sed -i '/cache:cache/i stats:github.com/GustavoKatel/coredns-stats' plugin.cfg

ADD . /coredns/plugin/stats

RUN echo "replace github.com/GustavoKatel/coredns-stats => /coredns/plugin/stats" >> go.mod
RUN go get github.com/GustavoKatel/coredns-stats

RUN go generate && make

FROM alpine:latest

COPY --from=builder /coredns/coredns /coredns/coredns

ENTRYPOINT ["/coredns/coredns"]

# Build image
FROM golang:1.15-alpine3.13 AS builder

ENV GOFLAGS="-mod=readonly"
ENV CGO_ENABLED=0

RUN apk add --update --no-cache bash ca-certificates curl git build-base

WORKDIR /usr/local/src/_build

ARG GOPROXY

ARG PLZ_BUILD_CONFIG
ARG PLZ_OVERRIDES
ARG PLZ_CONFIG_PROFILE

ENV PLZ_ARGS="-p -o \"build.path:${PATH}\""

COPY .plzconfig* pleasew ./
RUN ./pleasew update

COPY third_party third_party/
RUN ./pleasew build //third_party/...

COPY . .

RUN ./pleasew export outputs -o /usr/local/bin :blob-proxy


# Final image
FROM alpine:3.13.0

RUN apk add --update --no-cache ca-certificates tzdata bash curl libc6-compat

SHELL ["/bin/bash", "-c"]

# set up nsswitch.conf for Go's "netgo" implementation
# https://github.com/gliderlabs/docker-alpine/issues/367#issuecomment-424546457
RUN test ! -e /etc/nsswitch.conf && echo 'hosts: files dns' > /etc/nsswitch.conf

COPY --from=builder /usr/local/bin/ /usr/local/bin/

CMD blob-proxy --addr :${PORT:-8000}

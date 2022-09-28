# Simple usage with a mounted data directory:
# > docker build --platform="linux/amd64" -t simapp . --build-arg ARCH=x86_64
#
# Server:
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simapp:/root/.simapp simapp simd init test-chain
# TODO: need to set validator in genesis so start runs
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simapp:/root/.simapp simapp simd start
#
# Client: (Note the simapp binary always looks at ~/.simapp we can bind to different local storage)
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simappcli:/root/.simapp simapp simd keys add foo
# > docker run -it -p 26657:26657 -p 26656:26656 -v ~/.simappcli:/root/.simapp simapp simd keys list
# TODO: demo connecting rest-server (or is this in server now?)
FROM golang:alpine AS build-env
ARG ARCH=$ARCH

# Install minimum necessary dependencies,
ENV PACKAGES curl make git libc-dev bash gcc g++ linux-headers eudev-dev python3
RUN apk add --update --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/line/lbm-sdk

# prepare dbbackend before building; this can be cached
COPY ./Makefile ./
COPY ./contrib ./contrib
#RUN make dbbackend

# Install GO dependencies
COPY ./go.mod /go/src/github.com/line/lbm-sdk/go.mod
COPY ./go.sum /go/src/github.com/line/lbm-sdk/go.sum
RUN go mod download

# See https://github.com/line/wasmvm/releases
# See https://github.com/line/wasmvm/releases
ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.x86_64.a /lib/libwasmvm_static.x86_64.a
ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.aarch64.a /lib/libwasmvm_static.aarch64.a
RUN sha256sum /lib/libwasmvm_static.aarch64.a | grep bc3db72ba32f34ad88ceb1d20479411bd7f50ccd6a5ca50cc8ca462a561e6189
RUN sha256sum /lib/libwasmvm_static.x86_64.a | grep 352fa5de5f9dba66f0a38082541d3e63e21394fee3e577ea35e0906294c61276

RUN ln -s /lib/libwasmvm_static.${ARCH}.a /usr/lib/libwasmvm_static.a

# Add source files
COPY . .

# install simapp, remove packages
RUN BUILD_TAGS=static make build CGO_ENABLED=1

# Final image
FROM alpine:edge

# Set up OS dependencies
RUN apk add --update --no-cache  ca-certificates libstdc++
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/line/lbm-sdk/build/simd /usr/bin/simd

EXPOSE 26656 26657 1317 9090

# Run simd by default, omit entrypoint to ease using container with simcli
CMD ["simd"]

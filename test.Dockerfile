FROM golang:alpine AS build-env
ARG ARCH=$ARCH
ENV PACKAGES curl make git libc-dev bash gcc g++ linux-headers eudev-dev python3
RUN apk add --update --no-cache $PACKAGES

WORKDIR /go/src/github.com/line/lbm-sdk
COPY ./Makefile ./
COPY ./contrib ./contrib
COPY ./go.mod /go/src/github.com/line/lbm-sdk/go.mod
COPY ./go.sum /go/src/github.com/line/lbm-sdk/go.sum
RUN go mod download

ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.x86_64.a /lib/libwasmvm_static.x86_64.a
ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.aarch64.a /lib/libwasmvm_static.aarch64.a
RUN sha256sum /lib/libwasmvm_static.aarch64.a | grep bc3db72ba32f34ad88ceb1d20479411bd7f50ccd6a5ca50cc8ca462a561e6189
RUN sha256sum /lib/libwasmvm_static.x86_64.a | grep 352fa5de5f9dba66f0a38082541d3e63e21394fee3e577ea35e0906294c61276

RUN ln -s /lib/libwasmvm_static.${ARCH}.a /usr/lib/libwasmvm_static.a

COPY . .

RUN BUILD_TAGS=static make build CGO_ENABLED=1

FROM alpine:edge

# Set up OS dependencies
RUN apk add --update --no-cache  ca-certificates libstdc++ bash
WORKDIR /root
COPY --from=build-env /go/src/github.com/line/lbm-sdk/build/simd /usr/bin/simd
COPY --from=build-env /go/src/github.com/line/lbm-sdk/test_init_node.sh /root/init_node.sh
RUN chmod +x /root/init_node.sh && /root/init_node.sh sim 1
EXPOSE 26657 1317 9090 9091

# Run simd by default, omit entrypoint to ease using container with simcli
CMD ["simd", "start", "--home", "/root/.simapp/simapp0"]

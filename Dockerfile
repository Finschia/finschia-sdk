# Simple usage with a mounted data directory:
# > docker build --platform="linux/amd64" -t simapp .
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

# Build cosmwasm
ENV RUSTUP_HOME=/usr/local/rustup
ENV CARGO_HOME=/usr/local/cargo
ENV PATH=$CARGO_HOME/bin:$PATH

RUN wget "https://static.rust-lang.org/rustup/dist/x86_64-unknown-linux-musl/rustup-init"
RUN chmod +x rustup-init
RUN ./rustup-init -y --no-modify-path --default-toolchain 1.57.0; rm rustup-init
RUN chmod -R a+w $RUSTUP_HOME $CARGO_HOME
RUN cd $(go list -f "{{ .Dir }}" -m github.com/line/wasmvm) && \
    cd ./libwasmvm && \
    RUSTFLAGS='-C target-feature=-crt-static' cargo build --release --example staticlib && \
    mv -f target/release/examples/libstaticlib.a /usr/lib/libwasmvm_static.a && \
    rm -rf target

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

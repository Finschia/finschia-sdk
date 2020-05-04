# Simple usage with a mounted data directory:
# > docker build -t line/link .
# > docker run -it -p 26656:26656 -p 26657:26657 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli line/link linkd init
# > docker run -it -p 26656:26656 -p 26657:26657 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli line/link linkd start --rpc.laddr=tcp://0.0.0.0:26657 --p2p.laddr=tcp://0.0.0.0:26656
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python perl
ARG GITHUB_TOKEN=""
RUN apk add --no-cache $PACKAGES
# Set working directory for the build
WORKDIR /linkchain-build

COPY ./go.mod /linkchain-build/go.mod
COPY ./go.sum /linkchain-build/go.sum
RUN go env -w GOPRIVATE=github.com/line/*
# GITHUB_TOKEN should be provided to build link docker image
RUN git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN go mod download

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN  make install

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/bin/linkd /usr/bin/linkd
COPY --from=build-env /go/bin/linkcli /usr/bin/linkcli

# Run linkd by default, omit entrypoint to ease using container with linkcli
CMD ["linkd"]

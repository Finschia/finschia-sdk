# Simple usage with a mounted data directory:
# > docker build -t link .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli link linkd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli link linkd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python

# Set working directory for the build
WORKDIR /go/src/github.com/cosmos/link

# Add source files
COPY . .

# Install minimum necessary dependencies, build Cosmos SDK, remove packages
RUN apk add --no-cache $PACKAGES && \
    make tools && \
    make install

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

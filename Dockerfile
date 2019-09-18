# Simple usage with a mounted data directory:
# > docker build -t link .
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli link linkd init
# > docker run -it -p 46657:46657 -p 46656:46656 -v ~/.linkd:/root/.linkd -v ~/.linkcli:/root/.linkcli link linkd start
FROM golang:alpine AS build-env

# Set up dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /linkchain-build

COPY ./go.mod /linkchain-build/go.mod
COPY ./go.sum /linkchain-build/go.sum
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

ENTRYPOINT ["/usr/bin/wrapper.sh"]

# Run linkd by default, omit entrypoint to ease using container with linkcli
CMD ["start"]

COPY wrapper.sh /usr/bin/wrapper.sh
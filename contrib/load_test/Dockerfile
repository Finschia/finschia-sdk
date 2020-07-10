# Build Image
FROM golang:alpine AS build-env

WORKDIR /link-load-tester

ENV PACKAGES git
ARG GITHUB_TOKEN=""
RUN apk add --no-cache $PACKAGES

# Add source files
COPY . .

RUN ls

# Install dependencies
RUN git config --global url."https://$GITHUB_TOKEN:x-oauth-basic@github.com/".insteadOf "https://github.com/"
RUN go mod download

# Build modules
RUN GOOS=linux GOARCH=amd64 go build -o build/link-load-tester ./contrib/load_test/cmd/link_load_tester

# ---------------------------------
# Production Image
FROM alpine:edge

WORKDIR /link-load-tester

# Copy over binaries from the build-env
COPY --from=build-env /link-load-tester/build/link-load-tester /usr/bin/link-load-tester
COPY --from=build-env /link-load-tester/contrib/load_test/config.yaml /link-load-tester/contrib/load_test/config.yaml
EXPOSE 8000

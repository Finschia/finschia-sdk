FROM alpine:3.7
MAINTAINER LinkChain <dl_linkchain_v2_dev@linecorp.com>

RUN apk update && \
    apk upgrade && \
    apk add --update --no-cache autoconf bash jq file curl git lsof wget iputils vim busybox-extras

ENV LINKCLIHOME /linkd
ENV PATH="$PATH:$LINKCLIHOME"

STOPSIGNAL SIGTERM
COPY wrapper.sh /usr/bin/wrapper.sh
WORKDIR /linkd

EXPOSE 26656 26657
ENTRYPOINT ["/usr/bin/wrapper.sh"]
CMD ["start"]


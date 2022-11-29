FROM goreleaser/goreleaser-cross:v1.18.3

ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.x86_64.a /lib/libwasmvm_static.x86_64.a
ADD https://github.com/line/wasmvm/releases/download/v1.0.0-0.10.0/libwasmvm_static.aarch64.a /lib/libwasmvm_static.aarch64.a
RUN sha256sum /lib/libwasmvm_static.aarch64.a | grep bc3db72ba32f34ad88ceb1d20479411bd7f50ccd6a5ca50cc8ca462a561e6189
RUN sha256sum /lib/libwasmvm_static.x86_64.a | grep 352fa5de5f9dba66f0a38082541d3e63e21394fee3e577ea35e0906294c61276

RUN cp /lib/libwasmvm_static.x86_64.a /usr/lib/libwasmvm_static.a
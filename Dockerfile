ARG GOLANG_IMAGE=golang:1.17.6-alpine3.15

ARG BUILD_IMAGE=alpine:3.15

FROM $GOLANG_IMAGE as builder

COPY . /dist/core

RUN apk add \
		git \
		make && \
	cd /dist/core && \
	go version && \
	make release && \
	make import

FROM $BUILD_IMAGE

COPY --from=builder /dist/core/core /core/bin/core
COPY --from=builder /dist/core/import /core/bin/import
COPY --from=builder /dist/core/mime.types /core/mime.types
COPY --from=builder /dist/core/run.sh /core/bin/run.sh

RUN mkdir /core/config /core/data

ENV CORE_CONFIGFILE=/core/config/config.json
ENV CORE_STORAGE_DISK_DIR=/core/data
ENV CORE_DB_DIR=/core/config

VOLUME ["/core/data", "/core/config"]
ENTRYPOINT ["/core/bin/run.sh"]
WORKDIR /core

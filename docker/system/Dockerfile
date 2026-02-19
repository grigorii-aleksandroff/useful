ARG ALPINE_VERSION

FROM alpine:${ALPINE_VERSION}

CMD apk add --no-cache bash && \
    sysctl -w vm.max_map_count=262144 &&\
    tail -f /dev/null

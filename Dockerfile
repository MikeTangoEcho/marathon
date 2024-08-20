ARG FFMPEG_VERSION=latest
FROM mwader/static-ffmpeg:${FFMPEG_VERSION} as ffmpeg

FROM alpine:latest

# BEGIN Bins
# FFMPEG
COPY --from=ffmpeg /ffmpeg /usr/local/bin/

# MARATHON
COPY marathon /usr/local/bin/
# END Bins

RUN addgroup -g 1001 marathon \
    && adduser -D -H -u 1001 marathon -G marathon
USER marathon

ENTRYPOINT ["/usr/local/bin/marathon"]
CMD ["--help"]
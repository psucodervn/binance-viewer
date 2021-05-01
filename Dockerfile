ARG BINARY=copytrader
ARG DIR=/app

FROM alpine AS builder

RUN apk update && apk add --no-cache ca-certificates tzdata

FROM scratch
ARG BINARY

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY assets ./assets
COPY $BINARY ./app

ENTRYPOINT ["./app"]

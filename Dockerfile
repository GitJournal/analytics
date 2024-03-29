FROM golang:1.18 as builder

WORKDIR /app
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Deps
RUN apt-get update && apt-get install -y git-crypt

# Build the binary.
RUN go build -v -o server

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

ARG HEAD_SHA
ARG PGPASSWORD

ENV SENTRY_RELEASE=$HEAD_SHA

ARG MAXMIND_LICENSE_KEY
RUN curl -f -o GeoLite2-City.tar.gz "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=$MAXMIND_LICENSE_KEY&suffix=tar.gz" && \
    tar -xf GeoLite2-City.tar.gz && \
    mv GeoLite2-City_*/GeoLite2-City.mmdb / && \
    rm -rf GeoLite2-City.tar.gz GeoLite2-City
CMD ["/app/server"]

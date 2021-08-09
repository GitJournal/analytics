FROM golang:1.16-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server

FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates curl && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server

ARG MAXMIND_LICENSE_KEY
RUN curl -f -o ./GeoLite2-City.mmdb "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=$MAXMIND_LICENSE_KEY&suffix=tar.gz"
CMD ["/app/server"]

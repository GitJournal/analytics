MAXMIND=$(shell cat secrets/maxmind)
HEAD_SHA=$(shell git rev-parse --short HEAD)

build:
	docker build \
		--build-arg MAXMIND_LICENSE_KEY="$(MAXMIND)" \
		--build-arg HEAD_SHA="$(HEAD_SHA)" \
		-t "ghcr.io/gitjournal/analytics_backend:latest" .

push:
	docker push "ghcr.io/gitjournal/analytics_backend:latest"

.PHONY: protos

protos:
	buf generate

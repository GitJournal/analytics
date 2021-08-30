#protoc --go_out=. --go_opt=paths=source_relative \
#    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
#    protos/analytics.proto

MAXMIND=$(shell cat secrets/maxmind)
HEAD_SHA=$(shell git rev-parse --short HEAD)
PGPASSWORD=$(shell cat secrets/postgres)

build:
	docker build \
		--build-arg MAXMIND_LICENSE_KEY="$(MAXMIND)" \
		--build-arg HEAD_SHA="$(HEAD_SHA)" \
		--build-arg PGPASSWORD="$(PGPASSWORD)" \
		-t "gcr.io/gitjournal-io/analytics_backend" .

push:
	docker push gcr.io/gitjournal-io/analytics_backend:latest

deploy:
	gcloud run deploy analyticsbackend --image gcr.io/gitjournal-io/analytics_backend

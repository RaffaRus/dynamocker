# Get version in a <maj>.<minor> format
VERSION = $(shell git describe | sed 's/^v//')

# Docker image management

build-docker-image-fe:
		docker build -t dynamocker-fe:$(VERSION) -f docker/fe.Dockerfile --ulimit nofile=65535:65535  .

build-docker-image-be:
		docker build -t dynamocker-be:$(VERSION) -f docker/be.Dockerfile --ulimit nofile=65535:65535  .

build-docker-images: build-docker-image-fe build-docker-image-be

docker-compose-build-up:
		docker compose -f docker/docker-compose-build.yml --env-file docker/.env up -d

docker-compose-image-up:
		docker compose -f docker/docker-compose-image.yml --env-file docker/.env up -d
		
docker-compose-image-push:
		docker compose -f docker/docker-compose-image.yml --env-file docker/.env push
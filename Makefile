# Get version in a <maj>.<minor> format
VERSION = $(shell git describe | sed 's/^v//')


# Podman
build-podman-images: build-podman-image-fe build-podman-image-be

build-podman-image-fe:
		podman build -t dynamocker-fe:$(VERSION) -f docker/fe.Dockerfile --ulimit nofile=65535:65535  .
build-podman-image-be:
		podman build -t dynamocker-be:$(VERSION) -f docker/be.Dockerfile --ulimit nofile=65535:65535  .

run-podman-images: run-podman-image-fe run-podman-image-be

run-podman-image-be:
		podman run --name dynamocker-be -d -i localhost/dynamocker-be:$(VERSION) 
run-podman-image-fe:
		podman run --name dynamocker-fe -d -i localhost/dynamocker-fe:$(VERSION) 


# Docker

build-docker-image-fe:
		docker build -t dynamocker-fe:$(VERSION) -f docker/fe.Dockerfile --ulimit nofile=65535:65535  .

build-docker-image-be:
		docker build -t dynamocker-be:$(VERSION) -f docker/be.Dockerfile --ulimit nofile=65535:65535  .

build-docker-images: build-docker-image-fe build-docker-image-be

docker-run:
		docker run .
		
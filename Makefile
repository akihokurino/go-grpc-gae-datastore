MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export PATH := $(ROOT)/scripts:$(PATH)

GOOGLE_SERVICE_KEY := ""
SERVICE := default
VERSION := 1

vendor:
	go mod tidy

gen:
	mkdir -p proto/go
	rm -rf proto/go/*
	protoc --proto_path=proto/. --twirp_out=proto/go/ --go_out=proto/go/ proto/*.proto
	cp di/wire_gen.default.go di/wire_gen.go
	go generate di/wire_gen.go

build:
	GOOS=linux GOARCH=amd64 go build -o .tmp/default ./entrypoint/default/
	GOOS=linux GOARCH=amd64 go build -o .tmp/batch ./entrypoint/batch/
	GOOS=linux GOARCH=amd64 go build -o .tmp/subscriber ./entrypoint/subscriber/

encrypt:
	gcloud config set project akiho-playground
	sops --encrypt \
        --gcp-kms projects/akiho-playground/locations/asia-northeast1/keyRings/keys/cryptoKeys/config \
        config/env.raw.yaml > config/env.enc.yaml

decrypt:
	gcloud config set project akiho-playground
	sops --decrypt config/env.enc.yaml > config/env.raw.yaml

setup-config: decrypt
	setup_config.sh

deploy: build setup-config
	SERVICE=${SERVICE} VERSION=${VERSION} deploy.sh

deploy-all: build setup-config
	VERSION=${VERSION} COMMAND=all deploy.sh

deploy-function:
	deploy_function.sh

backup:
	backup.sh

run-local:
	docker-compose up

format:
	find . -print | grep --regex '.*\.go' | xargs goimports -w -local
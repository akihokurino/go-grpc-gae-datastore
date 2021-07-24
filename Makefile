MAKEFLAGS=--no-builtin-rules --no-builtin-variables --always-make
ROOT := $(realpath $(dir $(lastword $(MAKEFILE_LIST))))
export PATH := $(ROOT)/scripts:$(PATH)

GOOGLE_SERVICE_KEY := ""
SERVICE := default
VERSION := 1
ENV := dev

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
	gcloud config set project gae-go-sample-${ENV}
	sops --encrypt \
        --gcp-kms projects/gae-go-sample-${ENV}/locations/asia-northeast1/keyRings/gae-go-sample/cryptoKeys/config \
        config/${ENV}_env.raw.yaml > config/${ENV}_env.enc.yaml

decrypt:
	gcloud config set project gae-go-sample-${ENV}
	sops --decrypt config/${ENV}_env.enc.yaml > config/${ENV}_env.raw.yaml

setup-config: decrypt
	ENV=${ENV} setup_config.sh

deploy: build setup-config
	ENV=${ENV} SERVICE=${SERVICE} VERSION=${VERSION} deploy.sh

deploy-all: build setup-config
	ENV=${ENV} VERSION=${VERSION} COMMAND=all deploy.sh

deploy-function:
	ENV=${ENV} deploy_function.sh

backup:
	ENV=${ENV} backup.sh

run-local:
	docker-compose up

format:
	find . -print | grep --regex '.*\.go' | xargs goimports -w -local
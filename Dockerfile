FROM google/cloud-sdk:alpine AS ci-base
RUN apk update && apk add git curl build-base autoconf automake libtool make bash protobuf
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/bin:$PATH
ENV GO_VERSION 1.16.4
RUN curl -Lso go.tar.gz "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf go.tar.gz \
    && rm go.tar.gz
ENV PATH /usr/local/go/bin:$PATH
RUN GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    && GO111MODULE=off go get -u github.com/google/wire/cmd/wire \
    && go get go.mozilla.org/sops/v3/cmd/sops@v3.5.0 \
    && go get github.com/twitchtv/twirp/protoc-gen-twirp@v8.0.0 \
    && gcloud components install app-engine-go
ENV PATH /google-cloud-sdk/platform/google_appengine:$PATH


FROM ci-base AS cloud-build-base
ENV APP_ROOT /app
WORKDIR $APP_ROOT
COPY . $APP_ROOT/


FROM golang:1.16-alpine AS local-dev
RUN apk update && apk add --no-cache g++ gcc git make bash protobuf ca-certificates curl
RUN GO111MODULE=off go get -u google.golang.org/protobuf/cmd/protoc-gen-go \
    && GO111MODULE=off go get -u github.com/google/wire/cmd/wire \
    && go get github.com/twitchtv/twirp/protoc-gen-twirp@v8.0.0
ENV TZ=Asia/Tokyo
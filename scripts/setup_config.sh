#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..

go run ${APP_ROOT}/scripts/merge_yaml/main.go \
    ${APP_ROOT}/entrypoint/default/app.yaml \
    ${APP_ROOT}/entrypoint/default/app_template.yaml \
    ${APP_ROOT}/config/env.raw.yaml

go run ${APP_ROOT}/scripts/merge_yaml/main.go \
    ${APP_ROOT}/entrypoint/batch/app.yaml \
    ${APP_ROOT}/entrypoint/batch/app_template.yaml \
    ${APP_ROOT}/config/env.raw.yaml

go run ${APP_ROOT}/scripts/merge_yaml/main.go \
    ${APP_ROOT}/entrypoint/subscriber/app.yaml \
    ${APP_ROOT}/entrypoint/subscriber/app_template.yaml \
    ${APP_ROOT}/config/env.raw.yaml
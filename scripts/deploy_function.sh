#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..
PROJECT_ID=akiho-playground

cd ${APP_ROOT}/functions && npm install && cd ${APP_ROOT}

firebase use ${PROJECT_ID}
firebase deploy --only functions --project ${PROJECT_ID}
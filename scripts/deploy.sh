#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..

PROJECT_ID=akiho-playground

gcloud --quiet config set project ${PROJECT_ID}

deployConfig () {
    gcloud app deploy --quiet ${APP_ROOT}/entrypoint/default/cron.yaml &
    gcloud app deploy --quiet ${APP_ROOT}/entrypoint/default/index.yaml &
    wait
}

if [[ ${COMMAND} = "all" ]]; then
    cp ${APP_ROOT}/entrypoint/default/app.yaml ${APP_ROOT}/deploy_gae_default.yaml
    cp ${APP_ROOT}/entrypoint/batch/app.yaml ${APP_ROOT}/deploy_gae_batch.yaml
    cp ${APP_ROOT}/entrypoint/subscriber/app.yaml ${APP_ROOT}/deploy_gae_subscriber.yaml

    deployConfig

    gcloud app deploy --quiet --version ${VERSION} --project ${PROJECT_ID} \
        ${APP_ROOT}/deploy_gae_default.yaml
    gcloud app deploy --quiet --version ${VERSION} --project ${PROJECT_ID} \
        ${APP_ROOT}/deploy_gae_batch.yaml
    gcloud app deploy --quiet --version ${VERSION} --project ${PROJECT_ID} \
        ${APP_ROOT}/deploy_gae_subscriber.yaml

    rm ${APP_ROOT}/deploy_gae_default.yaml
    rm ${APP_ROOT}/deploy_gae_batch.yaml
    rm ${APP_ROOT}/deploy_gae_subscriber.yaml
else
    if [[ ${SERVICE} = "config" ]]; then
        deployConfig
    else
        cp ${APP_ROOT}/entrypoint/${SERVICE}/app.yaml ${APP_ROOT}/deploy_gae_${SERVICE}.yaml

        gcloud app deploy --quiet --version ${VERSION} --project ${PROJECT_ID} \
            ${APP_ROOT}/deploy_gae_${SERVICE}.yaml

        rm ${APP_ROOT}/deploy_gae_${SERVICE}.yaml
    fi
fi
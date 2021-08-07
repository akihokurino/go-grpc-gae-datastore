#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..

DATE=`date -u '+%Y%m%d'`
DATASET=DatastoreBackup
PROJECT_ID=akiho-playground
BUCKET=akiho-playground-datastore-backup

KINDS=ApplyClient,Client,Company,Contract,Customer,Entry,MessageRoom,Message,NoEntrySupport,NoMessageSupport,Project,User

gcloud config set project ${PROJECT_ID}

gcloud datastore export --namespaces="(default)" ${BUCKET}/${DATE}_all
gcloud datastore export --kinds=${KINDS} --namespaces="(default)" ${BUCKET}/${DATE}_each

bq --location=us-central1 mk -d --description "" ${DATASET}
kindList=(${KINDS//,/ })
for kind in ${kindList[@]}; do
  bq --location=us-central1 load --source_format=DATASTORE_BACKUP ${DATASET}.${kind}_${DATE} \
    gs://${BUCKET}/${DATE}_each/default_namespace/kind_${kind}/default_namespace_kind_${kind}.export_metadata
done

# 復元例:
# gcloud datastore import gs://akiho-playground-datastore-backup/20210523_all/20210523_all.overall_export_metadata
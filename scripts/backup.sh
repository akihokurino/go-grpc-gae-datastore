#!/usr/bin/env bash

APP_ROOT=$(dirname $0)/..

DATE=`date -u '+%Y%m%d'`
DATASET=DatastoreBackup
PROJECT_ID=gae-go-sample-${ENV}
BUCKET=gae-go-sample-${ENV}-datastore-backup

KINDS=ApplyClient,ClientBookmark,Client,Company,Contract,CustomerBilling,CustomerBlackList,CustomerMeta,Customer,Entry,ExternalContractActivity,ExternalContractBillingSchedule,ExternalContractBilling,ExternalContractInvoice,ExternalContract,JobCareer,MessageRoom,Message,NoEntrySupport,NoMessageSupport,Project,User

gcloud config set project ${PROJECT_ID}

gcloud datastore export --namespaces="(default)" ${BUCKET}/${DATE}_all
gcloud datastore export --kinds=${KINDS} --namespaces="(default)" ${BUCKET}/${DATE}_each

bq --location=asia-northeast1 mk -d --description "" ${DATASET}
kindList=(${KINDS//,/ })
for kind in ${kindList[@]}; do
  bq --location=asia-northeast1 load --source_format=DATASTORE_BACKUP ${DATASET}.${kind}_${DATE} \
    gs://${BUCKET}/${DATE}_each/default_namespace/kind_${kind}/default_namespace_kind_${kind}.export_metadata
done

# 復元例:
# gcloud datastore import gs://gae-go-sample-dev-datastore-backup/20210523_all/20210523_all.overall_export_metadata
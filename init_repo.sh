#!/bin/bash

set -e

CLIENT_ID='x'
CLIENT_SECRET='y'
ENDPOINT_URL='http://103.61.38.253:9100'
CONFIG='config.json'

mkdir -p ${HOME}/.openpitrix/ && cd ${HOME}/.openpitrix/
touch ${CONFIG}

echo -e "{" > ${CONFIG}
echo -e "\t\"client_id\": \"$CLIENT_ID\"," >> ${CONFIG}
echo -e "\t\"client_secret\": \"$CLIENT_SECRET\"," >> ${CONFIG}
echo -e "\t\"endpoint_url\": \"$ENDPOINT_URL\"" >> ${CONFIG}
echo -e "}" >> ${CONFIG}


tars=$(ls -l /data/helm-pkg |grep ".tgz" |awk '{print $NF}')
cd /data/helm-pkg
for pkg in $tars
do
    name=${pkg%-*}
    version_id=$(opctl create_app --version_package $pkg --name $name --version_type helm |grep \"version_id\": |awk -F ':' '{print $2}' |sed 's/\"//g')
    opctl submit_app_version --version_id ${version_id}
    opctl admin_pass_app_version --version_id ${version_id}
    opctl release_app_version --version_id ${version_id}
done

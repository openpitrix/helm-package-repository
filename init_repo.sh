#!/bin/bash

set -e

if [ ! -f "${HOME}/.openpitrix/config.json" ];then
  mkdir -p ${HOME}/.openpitrix/
  mv /usr/local/bin/config.json ${HOME}/.openpitrix/
fi

tars=$(ls -l /data/helm-pkg |grep ".tgz" |awk '{print $NF}')
cd /data/helm-pkg
for pkg in $tars
do
    name=${pkg%-*}
    version=${pkg##*-}
    version_name=${version%.*}
    version_id=$(opctl create_app --version_package $pkg --name $name --version_name $version_name --version_type helm |grep \"version_id\": |awk -F ':' '{print $2}' |sed 's/\"//g')
    opctl submit_app_version --version_id ${version_id}
    opctl admin_pass_app_version --version_id ${version_id}
    opctl release_app_version --version_id ${version_id}
done

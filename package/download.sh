#!/usr/bin/env bash

for d in $(cat ./config.json | jq '.[] | .appname');
do
    dir=$(echo $d | sed 's/\"//g')
    rm -rf $dir && mkdir -p $dir && cd $dir
    for url in $(cat ../config.json | jq --arg n $dir '.[] |select(.appname==$n)' | jq .versions[].url);
    do
        wget $(echo $url | sed 's/\"//g')
    done
    cd ..
done

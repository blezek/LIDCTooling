#!/bin/sh

cluster=$1

source venv/bin/activate
cluster_meta="$(starcluster listclusters $cluster 2>&1)"
# MASTER=$(echo "${cluster_meta}" | grep "master running" | sed -e 's/.*master.*\(ec2-.*com\)/\1/g')

MASTER=$(echo "${cluster_meta}" | grep "master running" | awk '{print $4}')

ssh-add  ~/.ssh/mykey.rsa
ssh-add  ~/.ssh/radiomics.rsa

rsync -avz sgeadmin@$MASTER:results.tar.gz .


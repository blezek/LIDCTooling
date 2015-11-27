#!/bin/sh

source venv/bin/activate
cluster_meta="$(starcluster listclusters $1 2>&1)"
# MASTER=$(echo "${cluster_meta}" | grep "master running" | sed -e 's/.*master.*\(ec2-.*com\)/\1/g')

MASTER=$(echo "${cluster_meta}" | grep "master running" | awk '{print $4}')

echo $MASTER

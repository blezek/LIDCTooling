#!/bin/sh

# SGE Parameters

## Name
#$ -N evaluate

## working directory
#$ -cwd

## Export variables
#$ -V

## Logs -- to /dev/null
# -o /home/sgeadmin/logs/
#$ -j yes
#$ -e /dev/null
#$ -o /dev/null

## Hold for lidc
#? -hold_jid lidc

cd /home/sgeadmin/

export PATH=$PATH:$JAVA_HOME/bin:/software/bin
export LD_LIBRARY_PATH=/software/lib

# Python Virtual env
source /software/lidc-venv/bin/activate

# Locations
SEGMENTED=/home/sgeadmin/segmented

DB=/home/sgeadmin/lidc.db

LIDCFetch --verbose \
    evaluate \
    --db $DB \
    $SEGMENTED/*






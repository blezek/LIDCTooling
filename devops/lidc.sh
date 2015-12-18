#!/bin/sh

# SGE Parameters

## Name
#$ -N lidc

## working directory
#$ -cwd

## Array job, should be 1-1318
#$ -t 1-1318

## Export variables
#$ -V

## Logs -- to /dev/null
# -o /home/sgeadmin/logs/
#$ -j yes
#$ -e /dev/null
#$ -o /dev/null

cd /home/sgeadmin/

export JAVA_HOME=/software/jvm/jdk1.8.0_20
export PATH=$PATH:$JAVA_HOME/bin:/software/bin
export LD_LIBRARY_PATH=/software/lib

# Give Java some room to work
export JAVA_OPTS="-Xmx2g"

# Python Virtual env
source /software/lidc-venv/bin/activate

# Substitute 1 if job is not set
export JOBID=${SGE_TASK_ID:=1}

# Get the XML file to run
XML=$(sed "${JOBID}q;d" /software/lidc.txt)

echo lidc.sh is processing $XML on `hostname --short`


# Paths of binaries
LESION=GenerateLesionSegmentation
EXTRACT=Extract
FETCH=LIDCFetch
EVALUATE=/software/python/evaluateSegmentation.py

# Locations
DICOM=/home/sgeadmin/dicom
SEGMENTED=/home/sgeadmin/segmented

DB=/home/sgeadmin/lidc.db

LIDCFetch --verbose \
    gather \
    --db $DB \
    --extract $EXTRACT \
    --fetch $FETCH \
    --lesion $LESION \
    --evaluate $EVALUATE \
    --dicom $DICOM \
    --segmented $SEGMENTED \
    --clean \
    /software/$XML






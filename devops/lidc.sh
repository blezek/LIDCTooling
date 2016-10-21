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

cd $HOME

export JAVA_HOME=/software/jvm/jdk1.8.0_20
export PATH=$PATH:$JAVA_HOME/bin:/software/bin:/software/python
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
FEATURES=/software/python/computeRadiomics.py
ALGORITHMS=/software/algorithms

# Local locations
DICOM=/tmp/dicom/$JOBID
SEGMENTED=/tmp/segmented/$JOBID

mkdir -p $DICOM $SEGMENTED

SeriesInstanceUID=`$EXTRACT SeriesInstanceUID /software/$XML`

if [ ! -e $HOME/segmented/$SeriesInstanceUID ]; then

    LIDCFetch  --verbose \
              process \
              --extract $EXTRACT \
              --fetch $FETCH \
              --evaluate $EVALUATE \
              --features $FEATURES \
              --dicom $DICOM \
              --segmented $SEGMENTED \
              --algorithms $ALGORITHMS \
              --clean-dicom \
              /software/$XML

    rsync -r $SEGMENTED/ $HOME/segmented

    rm -rf $SEGMENTED
    rm -rf $DICOM

else
    echo "Already processed $XML -- early exit because $HOME/segmented/$SeriesInstanceUID exists"
fi

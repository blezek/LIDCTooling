#!/bin/bash

# http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -eo pipefail
IFS=$'\n\t'

# Build instructions for Vagrant (Ubuntu trusty64)
NPROC=`getconf _NPROCESSORS_ONLN`

### Install our packages
sudo apt-get update 
sudo apt-get install -y cmake curl git subversion clang freeglut3-dev libxml2-dev g++ python-pip python-virtualenv python-dev
sudo apt-get install -y jq unzip
# freeglut3-dev brings in OpenGL

### Java
if [[ ! -e /usr/lib/jvm  ]]; then
   curl -L -O -H "Cookie: oraclelicense=accept-securebackup-cookie" -k "https://edelivery.oracle.com/otn-pub/java/jdk/8u20-b26/jdk-8u20-linux-x64.tar.gz"
   sudo mkdir -p /usr/lib/jvm
   sudo tar -C /usr/lib/jvm -xzf jdk-8u20-linux-x64.tar.gz
   sudo update-alternatives --install "/usr/bin/java" "java" "/usr/lib/jvm/jdk1.8.0_20/bin/java" 1000
   sudo update-alternatives --install "/usr/bin/javac" "javac" "/usr/lib/jvm/jdk1.8.0_20/bin/javac" 1000
   sudo update-alternatives --install "/usr/bin/javaws" "javaws" "/usr/lib/jvm/jdk1.8.0_20/bin/javaws" 1000
fi


### Chest Imaging Platform
cd $HOME
if [[ ! -e ChestImagingPlatform  ]]; then
    git clone https://github.com/acil-bwh/ChestImagingPlatform.git
    (cd ChestImagingPlatform && git checkout develop)
fi
### Don't rebuild if we don't have to...
if [[ ! -e /vagrant/ClusterSoftware/bin/GenerateLesionSegmentation ]]; then
    mkdir -p ChestImagingPlatform-build
    cd ChestImagingPlatform-build
    cmake ../ChestImagingPlatform/
    make -j $NPROC
    make
fi

### Build the LIDC code
cd
rsync --exclude ClusterSoftware --exclude segmented --exclude dicom --exclude '*journal' -ra /vagrant/ LIDCTooling

cd LIDCTooling
make build
./gradlew installDist

# jq and SQLite3 are not installed on the StarCluster machines
# so build and install
cd
if [[ ! -e jq ]]; then
    wget https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64
    mv jq-linux64 jq
    chmod 755 jq
fi

cd
if [[ ! -e sqlite3 ]]; then
    wget https://sqlite.org/2016/sqlite-amalgamation-3150000.zip
    unzip sqlite-amalgamation-3150000.zip
    cd sqlite-amalgamation-3150000
    cc -o sqlite3 sqlite3.c shell.c -lpthread -ldl
    cp sqlite3 ../
    cd ..
    rm -rf sqlite-amalgamation-3150000*
fi

### Build the python virtual environment
cd
virtualenv lidc-venv
source lidc-venv/bin/activate
pip install numpy==1.11.0
pip install -f http://www.simpleitk.org/SimpleITK/resources/software.html SimpleITK==0.9.1
pip install tqdm==4.7.1
pip install PyWavelets==0.4.0
(cd /pyradiomics && python setup.py install)

### Get LIDC XML files
cd
if [[ ! -e /vagrant/Clustersoftware/tcia-lidc-xml ]]; then
    wget -O LIDC_XML-only.tar.gz "https://wiki.cancerimagingarchive.net/download/attachments/3539039/LIDC-XML-only.tar.gz?version=1&modificationDate=1360694838194&api=v2"
    tar fxz LIDC_XML-only.tar.gz
    # LIDC XML files
    rsync -ar tcia-lidc-xml /vagrant/ClusterSoftware/
    find tcia-lidc-xml -name "*.xml" | sort > /vagrant/ClusterSoftware/lidc.txt
fi

### Copy to host
cd
# rm -rf /vagrant/ClusterSoftware/
mkdir -p /vagrant/ClusterSoftware/{bin,lib}
rsync LIDCTooling/bin/* /vagrant/ClusterSoftware/bin
rsync -ra LIDCTooling/build/install/LIDCTooling/bin/ /vagrant/ClusterSoftware/bin/
rsync -ra LIDCTooling/build/install/LIDCTooling/lib/ /vagrant/ClusterSoftware/lib/

# If we have built the CIP, install it
if [[ -e ChestImagingPlatform-build/ITKv4-build/lib/ ]]; then
    # install ITK libs
    rsync -ra --exclude "*.a" ChestImagingPlatform-build/ITKv4-build/lib/ /vagrant/ClusterSoftware/lib/

    # bin
    rsync -ra ChestImagingPlatform-build/CIP-build/bin/ /vagrant/ClusterSoftware/bin/
fi

# Python
rsync -ra lidc-venv /vagrant/ClusterSoftware
rsync -ra LIDCTooling/python /vagrant/ClusterSoftware/
rsync -ra LIDCTooling/algorithms /vagrant/ClusterSoftware/
sed -i.bak  s^/home/vagrant^/software^g /vagrant/ClusterSoftware/lidc-venv/bin/activate

# Java
rsync -ar /usr/lib/jvm /vagrant/ClusterSoftware/

# jq
rsync -ar jq /vagrant/ClusterSoftware/bin

# sqlite3
rsync -ar sqlite3 /vagrant/ClusterSoftware/bin

# Install everything locally
sudo ln --force --symbolic /vagrant/ClusterSoftware /software

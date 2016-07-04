#!/bin/sh

# Build instructions for Vagrant (Ubuntu trusty64)
NPROC=`getconf _NPROCESSORS_ONLN`

### Install our packages
sudo apt-get update 
sudo apt-get install -y cmake curl git subversion clang freeglut3-dev libxml2-dev g++ python-pip python-virtualenv
# freeglut3-dev brings in OpenGL

### Go
if [[ ! -e /usr/local/go/bin/go  ]]; then
   curl -L -O "https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz"
   sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
   cat >> .bashrc <<EOF
export PATH=$PATH:/usr/local/go/bin
EOF
   . .bashrc
fi

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
if [[ ! -e ChestImagingPlatform-build/CIP-build/bin/GenerateLesionSegmentation ]]; then
    mkdir -p ChestImagingPlatform-build
    cd ChestImagingPlatform-build
    cmake ../ChestImagingPlatform/
    make -j $NPROC
    make
fi

### Build the LIDC code
cd
branch=$(cd /vagrant && git branch --no-color 2> /dev/null | sed -e '/^[^*]/d' -e 's/* \(.*\)/\1/')
if [[ ! -e LIDCTooling ]]; then
   git clone https://github.com/dblezek/LIDCTooling.git
fi

# rsync --exclude ClusterSoftware --exclude segmented --exclude dicom -ra /vagrant LIDCTooling 

cd LIDCTooling
git pull
git checkout $branch
make build
./gradlew installDist


### Build the python virtual environment
cd
virtualenv lidc-venv
source lidc-venv/bin/activate
pip install -U pip
easy_install -f http://www.simpleitk.org/SimpleITK/resources/software.html SimpleITK

### Get LIDC XML files
cd
wget -O LIDC_XML-only.tar.gz "https://wiki.cancerimagingarchive.net/download/attachments/3539039/LIDC-XML-only.tar.gz?version=1&modificationDate=1360694838194&api=v2"
tar fxz LIDC_XML-only.tar.gz

### Copy to host
cd
rm -rf /vagrant/ClusterSoftware/
mkdir -p /vagrant/ClusterSoftware/{bin,lib}
rsync LIDCTooling/bin/* /vagrant/ClusterSoftware/bin
rsync -ra LIDCTooling/build/install/LIDCTooling/bin/ /vagrant/ClusterSoftware/bin/
rsync -ra LIDCTooling/build/install/LIDCTooling/lib/ /vagrant/ClusterSoftware/lib/

# install ITK libs
rsync -ra --exclude "*.a" ChestImagingPlatform-build/ITKv4-build/lib/ /vagrant/ClusterSoftware/lib/

# bin
rsync -ra ChestImagingPlatform-build/CIP-build/bin/ /vagrant/ClusterSoftware/bin/

# Python
rsync -ra lidc-venv /vagrant/ClusterSoftware
rsync -ra LIDCTooling/python /vagrant/ClusterSoftware/
rsync -ra LIDCTooling/algorithms /vagrant/ClusterSoftware/
sed -i.bak  s^/home/vagrant^/software^g /vagrant/ClusterSoftware/lidc-venv/bin/activate

# Java
rsync -ar /usr/lib/jvm /vagrant/ClusterSoftware/

# LIDC XML files
rsync -ar tcia-lidc-xml /vagrant/ClusterSoftware/
find tcia-lidc-xml -name "*.xml" | sort > /vagrant/ClusterSoftware/lidc.txt



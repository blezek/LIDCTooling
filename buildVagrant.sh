#!/bin/sh

# Build instructions for Vagrant (Ubuntu trusty64)

### Install our packages
sudo apt-get install -y cmake curl git subversion clang freeglut3-dev libxml2-dev g++
# freeglut3-dev brings in OpenGL

### Go
curl -L -O "https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz"
sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz
cat >> .bashrc <<EOF
export PATH=$PATH:/usr/local/go/bin
EOF
. .bashrc

### Java
curl -L -O -H "Cookie: oraclelicense=accept-securebackup-cookie" -k "https://edelivery.oracle.com/otn-pub/java/jdk/8u20-b26/jdk-8u20-linux-x64.tar.gz"
sudo mkdir -p /usr/lib/jvm
sudo tar -C /usr/lib/jvm -xzf jdk-8u20-linux-x64.tar.gz

sudo update-alternatives --install "/usr/bin/java" "java" "/usr/lib/jvm/jdk1.8.0_20/bin/java" 1000
sudo update-alternatives --install "/usr/bin/javac" "javac" "/usr/lib/jvm/jdk1.8.0_20/bin/javac" 1000
sudo update-alternatives --install "/usr/bin/javaws" "javaws" "/usr/lib/jvm/jdk1.8.0_20/bin/javaws" 1000

### Chest Imaging Platform
git clone https://github.com/acil-bwh/ChestImagingPlatform.git
mkdir ChestImagingPlatform-build
cd ChestImagingPlatform-build
cmake ../ChestImagingPlatform/
make -j 4
make

### Build the LIDC code
cd


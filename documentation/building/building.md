# Building the tools

# Chest Imaging Platform

Building the [Chest Imaging Platform](https://github.com/acil-bwh/ChestImagingPlatform) is relatively easy:

```bash
git clone https://github.com/acil-bwh/ChestImagingPlatform.git
cd ChestImagingPlatform
mkdir build
cd build
cmake ..
make # wait for it...
```

## LIDCFetch

`LIDCFetch` interacts with the LIDC website to fetch DICOM images.

```
git clone git@github.com:dblezek/LIDCTooling.git
cd LIDCTooling
# Build GO
make build
```

## Extract

The Extract tool takes DICOM images and a LIDC XML file and builds output images.

```
# Build Java app for local use
./gradlew installDist
# Build Java app for distribution
./gradlew distZip distTar
```

The built application can be run as:

```
build/install/LIDCTooling/bin/Extract
```

# Building using Vagrant

[Vagrant](https://www.vagrantup.com/) is useful for pre-compiling the tools rather than compiling on the (expensive) cluster.  Bring up a Vagrant image, compile all the source and install back in the host computer in the `ClusterSoftware` directory.

```
vagrant up
vagrant ssh
. /vagrant/buildVagrant.sh
```



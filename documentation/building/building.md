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

# LIDCTooling

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





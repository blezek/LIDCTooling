---
title: "Processing Overview"
weight: -100
---

# LIDC Data conversion

Making LIDC annotations available as DICOM Segmentation Objects.

## Abstract

The Lung Imaging Database Consortium has collected 1000's of CT scans with hand annotated tumor outlines for nodules > 3 mm.  Nodules smaller than 3 mm are indicated by a single "center of mass" annotation.  This data is stored in XML format and includes 4 readers.  Unfortunately, the format does not lend itself to easy comparison nor analysis.  This project aims to convert LIDC data to a more useful format.

## Goal

- convert LIDC XML data to DICOM segmentation objects
- export NIfTI volumes of segmentation masks

## Approach

### Automation

This process shall be automated to the extent possible, creating an easy to use, repeatable process for researchers desiring to obtain the LIDC data in a usable format.

Automation includes:

- download of CT datasets given a LIDC XML
- interaction with the [LIDC REST services](https://wiki.cancerimagingarchive.net/display/Public/TCIA+Programmatic+Interface+%28REST+API%29+Usage+Guide) ([PDF](TCIA REST API.pdf))

### Language / Platform

While [Slicer](http://www.slicer.org) is a natural choice, the process is not likely to be repeated nor does it require integration with Slicer.  LIDC data consists of XML and DICOM, and Java has easy to use libraries for both formats.  Java is natively available on every platform and integrates well with REST services.


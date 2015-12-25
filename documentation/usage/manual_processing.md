# Processing

## Example

### Get the LIDC XML Files

```
wget 'https://wiki.cancerimagingarchive.net/download/attachments/3539039/LIDC-XML-only.tar.gz'
mkdir -p LIDC-XML-only
cd LIDC-XML-only
tar fxvz ../LIDC-XML-only.tar.gz
cd ..
```

### Set some useful variables

```
XML=LIDC-XML-only/tcia-lidc-xml/157/158.xml
APIKEY=25f0025c-071c-426d-b15a-199421e2e889
APIKEY=864dcc73-ce40-4f19-8a3e-fce71fc2dba2
```

### Extract a SeriesInstanceUID from an XML file

```
SeriesInstanceUID=$(build/install/LIDCTooling/bin/Extract SeriesInstanceUID $XML)

# 1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824
echo $SeriesInstanceUID
```

### Download the image data

```
mkdir -p dicom/$SeriesInstanceUID
wget -O /tmp/$SeriesInstanceUID.zip --quiet --header "api_key: $APIKEY" "https://services.cancerimagingarchive.net/services/v3/TCIA/query/getImage?SeriesInstanceUID=$SeriesInstanceUID"

unzip -qq -o /tmp/$SeriesInstanceUID.zip -d dicom/$SeriesInstanceUID
```

### Extract the ROIs and generate JSON

```
mkdir -p segmented/$SeriesInstanceUID
build/install/LIDCTooling/bin/Extract segment $XML dicom/$SeriesInstanceUID segmented/$SeriesInstanceUID
```

### `reads.json`

The final result of the processing is a file named [segmented/$SeriesInstanceUID/reads.json](usage/reads.json).  It contains information about the original DICOM images, each read, nodules and small nodules found by each radiologist and pointers to a label map for each.


## Procesing in bulk

The `LIDCFetch` application can automatically process an entire `XML` file containing reads.

Example:

```bash
make build
bin/LIDCFetch gather LIDC-XML-only/tcia-lidc-xml/157/158.xml
```

Example of processing all the data:

```bash
make build
find LIDC-XML-only/tcia-lidc-xml -name "*.xml" | sort | xargs bin/LIDCFetch --verbose gather
# First 20
find LIDC-XML-only/tcia-lidc-xml -name "*.xml" | sort > all_xml_files.txt
head -20 all_xml_files.txt |  xargs bin/LIDCFetch --verbose gather --db Evaluation.db
```

Then collect the results into an SQLite3 database

```
find segmented | xargs bin/LIDCFetch evaluate --db Evaluation.db
```

## Tools

### Download DICOM datasets

Downloading DICOM from LIDC is relatively easy.  With an API key issued by the LIDC support team, issue `curl` commands with an extra header of `api_key`, e.g.

```
curl --header "api_key: 25f0025c-071c-426d-b15a-199421e2e889" https://services.cancerimagingarchive.net/services/v3/TCIA/query/getCollectionValues
```

### Download API

Usage:

```
LIDCFetch <command> [options]

Commands:

  collections -- get collections
  
  patients    -- get list of patient
    --collection <Collection>  -- collection query value
```

### XML To NIfTI

Load an XML file and convert to NIfTI.


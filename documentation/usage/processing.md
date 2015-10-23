# Processing

## Example



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


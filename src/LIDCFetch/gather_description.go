package main

var GatherDescription = `The gather command processes a LIDC read by performing the following steps:

1.  interrogates an XML file for imaging information
2.  downloads the image data to the "dicom/SeriesInstanceUID" directory
3.  creates a "reads.json" file in "segmentation/SeriesInstanceUID"
4.  creates NIfTi images in "segmented/SeriesInstanceUID"
    - each NIfTi file follows the naming convention "read_{#}.nii.gz"
      - "#" is the read number
5.  runs the "GenerateLesionSegmentation" algorithm on each read-nodule pair
    - output files follow the naming convention "read_{#}_nodule_{NID}.nii.gz"
      - "#" is the read number
      -  "NID" is the normalized nodule id
6.  evaluates the performance of the segmentation algorithm using "python/evaluateSegmentation.py"
    - writes a JSON file following the naming convention "read_{#}_nodule_{NID}_eval.json"
      - "#" and "NID" are as above

Location of the binaries are controlled by the following flags:

  --dicom       -- location for DICOM data
  --segmented   -- location for segmented files
  --extract     -- path to Extract Java program
  --fetch       -- path to LIDCFetch program
  --lesion      -- path to GenerateLesionSegmentation
  --evaluate    -- path to evaluateSegmentation.py

  NB: these locations have reasonable defaults.

`

# Evaluate lung nodule segmentation

usage = """
Usage:
 evaluateSegmentation.py [options] segmentation.nii gold_standard.nii evaluation.json
 options:
    --label <value>      -- label in gold_standard to compare, default is 1
    --threshold <value>  -- threshold in segmentation.nii, default is -0.5
"""

import getopt
import SimpleITK as sitk
import sys,os,json


# Parse 
opts,args = getopt.getopt(sys.argv[1:], "", longopts=["label=","threshold="])

if len(args) < 3:
    print usage
    sys.exit(1)

settings = { "--threshold": '-0.5', "--label": '1'}
settings.update ( opts )

# Load the input image
segmentation = sitk.ReadImage ( args[0] )
gold_standard = sitk.ReadImage ( args[1] )
jsonOutput = args[2]

segmentation = sitk.BinaryThreshold(segmentation, lowerThreshold=float(settings["--threshold"]), upperThreshold=10e10)
label = float(settings['--label'])
gold_standard = sitk.BinaryThreshold(gold_standard, lowerThreshold=label, upperThreshold=label)

# Compute overlap
ruler = sitk.LabelOverlapMeasuresImageFilter()
overlap = ruler.Execute(segmentation,gold_standard)

measures = {
    "false_negative_error": ruler.GetFalseNegativeError(),
    "false_positive_error": ruler.GetFalsePositiveError(),
    "mean_overlap": ruler.GetMeanOverlap(),
    "union_overlap": ruler.GetUnionOverlap(),
    "volume_similarity": ruler.GetVolumeSimilarity(),
    "jaccard_coefficient": ruler.GetJaccardCoefficient(),
    "dice_coefficient": ruler.GetDiceCoefficient()
    }

fid = open(jsonOutput, 'w')
fid.write(json.dumps(measures, indent=2))
fid.close()

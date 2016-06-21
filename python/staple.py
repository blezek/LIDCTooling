# Evaluate lung nodule segmentation

usage = """
Usage:
 staple.py [options] input0 [input1, ..., inputN]
 options:
    --label <value>        -- label of nodule
    --ignore <file>        -- skip this file
    --output <filename>    -- output filename
"""

import getopt
import SimpleITK as sitk
import sys,os,json


# Parse 
opts,args = getopt.getopt(sys.argv[1:], "", longopts=["label=","ignore=","output="])

if len(args) < 1:
    print usage
    sys.exit(1)

settings = { "--ignore": '', "--label": '1', '--output': ''}
settings.update ( opts )

# Load the input images
inputs = []
for f in args:
    if settings['--ignore'] != '' and f.find(settings['--ignore']) != -1:
        continue
    inputs.append ( sitk.ReadImage(f) )
    
staple = sitk.STAPLEImageFilter()
staple.SetForegroundValue ( float(settings["--label"]) )
output = staple.Execute ( *inputs )

sitk.WriteImage(output, settings['--output'], True)

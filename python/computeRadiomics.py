#!/usr/bin/env python
# Compute pyradiomics statistics on an image and a segmentation
# Output is saved as JSON

from radiomics import firstorder, glcm, imageoperations, shape, rlgl, glszm
import SimpleITK as sitk
import sys, os
import getopt,json
import logging


usage = """
computeRadiomics.py calculates radiomics features from an image and a mask.

Processing using the pyradiomics library from
https://github.com/Radiomics/pyradiomics Images must be readable
by SimpleITK (http://www.simpleitk.org/) and output is in JSON.

Usage:
  computeRadiomics.py [options] image mask output
  options:
    --log <value>   -- log level, DEBUG,INFO,WARNING default is WARNING
  parameters:
    image           -- image from which to compute statistics
    mask            -- image containing binary segmentation
    output          -- JSON output file
"""

#testBinWidth = 25 this is the default bin size
#testResampledPixelSpacing = (3,3,3) no resampling for now.

in_opts,args = getopt.getopt(sys.argv[1:], "", longopts=["label=","log="])
opts = { '--log': 'WARNING' }
opts.update ( in_opts )

if len(args) < 3:
  print ('\nNot enough arguments\n')
  print usage
  sys.exit(1)

imageName = args[0]
maskName = args[1]
outputName = args[2]
features = {}

level = logging.WARNING
if '--log' in opts:
  level = getattr(logging, opts['--log'].upper())
logging.basicConfig(level=level)

if not os.path.exists(imageName):
  logging.warning( 'Error: problem finding input image ' + imageName)
  os.exit(1)
if not os.path.exists(maskName):
  logging.warning ('Error: problem finding input image ' + maskName)
  os.exit(1)

image = sitk.ReadImage(imageName)
mask = sitk.ReadImage(maskName)

mask = sitk.BinaryThreshold(mask,1, 10000, 1, 0)
# sitk.WriteImage(mask, 'mask.nii')

apporx, ret = imageoperations.swt3(image)
image = ret[0].values()[0]
#
# Show the first order feature calculations
#

firstOrderFeatures = firstorder.RadiomicsFirstOrder(image,mask)

firstOrderFeatures.enableFeatureByName('Mean', True)
# firstOrderFeatures.enableAllFeatures()

# print 'Will calculate the following first order features: '
# for f in firstOrderFeatures.enabledFeatures.keys():
  # print '  ',f
  # print eval('firstOrderFeatures.get'+f+'FeatureValue.__doc__')

logging.info( 'Calculating first order features...')
firstOrderFeatures.calculateFeatures()
logging.info( 'done')

fof = {}
logging.info( 'Calculated first order features: ')
for (key,val) in firstOrderFeatures.featureValues.iteritems():
  logging.info( '  ' + key + ':' + str(val))
  fof[key] = val

  # now calculate the LoG features, following
  # steps and parameters as specified by @vnarayan13 in #22
  # 1. Get the LoG filtered image
  mmif = sitk.MinimumMaximumImageFilter()
  mmif.Execute(image)
  lowerThreshold = 0
  upperThreshold = mmif.GetMaximum()

  threshImage = imageoperations.applyThreshold(image,lowerThreshold=lowerThreshold, upperThreshold=upperThreshold,outsideValue=0)
  # get the mask of the thresholded pixels
  threshImageMask = imageoperations.applyThreshold(image,lowerThreshold=lowerThreshold, upperThreshold=upperThreshold,outsideValue=0,insideValue=1)
  # only include the voxels that are within the threshold
  threshMask = sitk.Cast(mask,1) & sitk.Cast(threshImageMask,1)
  import numpy
  sigmaValues = numpy.arange(5.,0.,-.5)[::1]
  for sigma in sigmaValues:
    logImage = imageoperations.applyLoG(image,sigmaValue=sigma)
    logFirstorderFeatures = firstorder.RadiomicsFirstOrder(logImage,threshMask)
    logFirstorderFeatures.enableAllFeatures()
    logFirstorderFeatures.calculateFeatures()
    logging.info( 'Calculated firstorder features with LoG sigma ' + str(sigma))
    for (key,val) in logFirstorderFeatures.featureValues.iteritems():
      laplacianFeatureName = 'LoG_sigma_%s_%s' %(str(sigma),key)
      logging.info( '  ' + laplacianFeatureName + ':' + str(val))
      fof[laplacianFeatureName] = val

features['first_order'] = fof

#
# Show Shape features
#
shapeFeatures = shape.RadiomicsShape(image, mask)
shapeFeatures.enableAllFeatures()

# logging.info( 'Will calculate the following Shape features: ')
# for f in shapeFeatures.enabledFeatures.keys():
  # print '  ',f
  # print eval('shapeFeatures.get'+f+'FeatureValue.__doc__')

logging.info( 'Calculating Shape features...')
shapeFeatures.calculateFeatures()
logging.info( 'done')

features['shape'] = {}
logging.info( 'Calculated Shape features: ')
for (key,val) in shapeFeatures.featureValues.iteritems():
  logging.info( '  ' + key + ':' + str(val))
  features['shape'][key] = val


#
# Show GLCM features
#
glcmFeatures = glcm.RadiomicsGLCM(image, mask, binWidth=25)
glcmFeatures.enableAllFeatures()

logging.info( 'Will calculate the following GLCM features: ')
# for f in glcmFeatures.enabledFeatures.keys():
  # print '  ',f
  # print eval('glcmFeatures.get'+f+'FeatureValue.__doc__')

logging.info( 'Calculating GLCM features...',)
glcmFeatures.calculateFeatures()
logging.info( 'done')

features['glcm'] = {}
logging.info( 'Calculated GLCM features: ')
for (key,val) in glcmFeatures.featureValues.iteritems():
  logging.info( '  '+key+':'+str(val))
  features['glcm'][key] = val

#
# Show RLGL features
#
rlglFeatures = rlgl.RadiomicsRLGL(image, mask, binWidth=25)
rlglFeatures.enableAllFeatures()

# logging.info( 'Will calculate the following RLGL features: ')
# for f in rlglFeatures.enabledFeatures.keys():
#   print '  ',f
#   print eval('rlglFeatures.get'+f+'FeatureValue.__doc__')

logging.info( 'Calculating RLGL features...')
rlglFeatures.calculateFeatures()
logging.info( 'done')

features['rlgl'] = {}
logging.info( 'Calculated RLGL features: ')
for (key,val) in rlglFeatures.featureValues.iteritems():
  logging.info( '  ' + key + ':' + str(val))
  features['rlgl'][key] = val

#
# Show GLSZM features
#
glszmFeatures = glszm.RadiomicsGLSZM(image, mask, binWidth=25)
glszmFeatures.enableAllFeatures()

# logging.info( 'Will calculate the following GLSZM features: ')
# for f in glszmFeatures.enabledFeatures.keys():
#   print '  ',f
#   print eval('glszmFeatures.get'+f+'FeatureValue.__doc__')

logging.info( 'Calculating GLSZM features...')
glszmFeatures.calculateFeatures()

features['glszm'] = {}
logging.info( 'Calculated GLSZM features: ')
for (key,val) in glszmFeatures.featureValues.iteritems():
  logging.info( '  '+key+':'+str(val))
  features['glszm'][key] = val

# Save
with open(outputName, 'w') as fid:
  logging.info( 'writing JSON' )
  json.dump(features,fid)

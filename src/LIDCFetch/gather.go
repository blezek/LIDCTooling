package main

import (
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GatherCommand = cli.Command{
	Name:   "gather",
	Action: gather,
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func Run(name string, arg ...string) (string, error) {
	t, err := exec.Command(name, arg...).CombinedOutput()
	return strings.TrimSpace(string(t)), err
}

func gather(context *cli.Context) {
	// Arguments are a list of XML files
	for _, xml := range context.Args() {
		SeriesInstanceUID, _ := Run("build/install/LIDCTooling/bin/Extract", "SeriesInstanceUID", xml)
		logger.Info("Processing %v - %v", xml, SeriesInstanceUID)
		// Download needed?
		DownloadDir := filepath.Join("dicom", SeriesInstanceUID)
		if !Exists(DownloadDir) {
			os.MkdirAll(DownloadDir, os.ModePerm|os.ModeDir)
			logger.Info("will be downloading to %v", DownloadDir)
		}
		SegmentedDir := filepath.Join("segmented", SeriesInstanceUID)
		logger.Info("%v exists? %v", SegmentedDir, Exists(SegmentedDir))
		if !Exists(SegmentedDir) {
			logger.Info("Segmenting %v", SegmentedDir)
			os.MkdirAll(SegmentedDir, os.ModePerm|os.ModeDir)
			Run("build/install/LIDCTooling/bin/Extract", "segment", xml, DownloadDir, SegmentedDir)
		}

		// Now parse the JSON
		JsonFile := filepath.Join(SegmentedDir, "reads.json")
		fid, _ := os.Open(JsonFile)

		everything, _ := jason.NewObjectFromReader(fid)
		originalFile, _ := everything.GetString("filename")
		reads, _ := everything.GetObjectArray("reads")
		for _, read := range reads {
			// Looping over all the reads
			logger.Debug("read: %+v", read)
			filename, _ := read.GetString("filename")
			readId, _ := read.GetInt64("id")
			logger.Debug("Filename: %v", filename)
			nodules, _ := read.GetObjectArray("nodules")
			for _, nodule := range nodules {
				id, _ := nodule.GetString("id")
				centroid, _ := nodule.GetFloat64Array("centroidLPS")
				label_value, _ := nodule.GetFloat64("label_value")
				logger.Debug("Nodule id: %v label: %v centroid: %v", id, label_value, centroid)

				// Process the lesion
				GLS := "/Users/blezek/Source/ChestImagingPlatform/build/CIP-build/bin/GenerateLesionSegmentation"
				input := filepath.Join(SegmentedDir, originalFile)
				output := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v.nii.gz", readId, id))
				seed := fmt.Sprintf("%v,%v,%v", centroid[0], centroid[1], centroid[2])
				out, err := Run(GLS, "-i", input, "-o", output, "--seeds", seed)
				logger.Info("Command line: %v", fmt.Sprintf("%v -i %v -o %v --seeds %v", GLS, input, output, seed))
				if err != nil {
					logger.Error("Error running %v -- %v\nOutput:%v", GLS, err.Error(), out)
					logger.Error("Command line: %v", fmt.Sprintf("%v -i %v -o %v --seeds", GLS, input, output, seed))
				}
			}
		}
	}
}

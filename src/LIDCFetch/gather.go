package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
)

var GatherCommand = cli.Command{
	Name:        "gather",
	Usage:       "<XML>",
	Description: GatherDescription,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "extract",
			Value: "build/install/LIDCTooling/bin/Extract",
			Usage: "Path to Extract Java program of LIDCTooling",
		},
		cli.StringFlag{
			Name:  "fetch",
			Value: "bin/LIDCFetch",
			Usage: "Path to LIDCFetch program of LIDCTooling",
		},
		cli.StringFlag{
			Name:  "lesion",
			Value: "/Users/blezek/Source/ChestImagingPlatform/build/CIP-build/bin/GenerateLesionSegmentation",
			Usage: "Path to GenerateLesionSegmentation program of Chest Imaging Toolkit",
		},
		cli.StringFlag{
			Name:  "evaluate",
			Value: "python/evaluateSegmentation.py",
			Usage: "Path to evaluateSegmentation.py program of LIDCTooling",
		},
		cli.StringFlag{
			Name:  "dicom",
			Value: "dicom",
			Usage: "location for DICOM data downloaded from LIDC",
		},
		cli.StringFlag{
			Name:  "segmented",
			Value: "segmented",
			Usage: "location for segmented files",
		},
		cli.BoolFlag{
			Name:  "clean",
			Usage: "Clean up all image files, both DICOM and NIfTI",
		},
		cli.BoolFlag{
			Name:  "clean-dicom",
			Usage: "Clean up all DICOM files",
		},
	},
	Action: gather,
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func Run(arg ...string) (string, error) {
	t, err := exec.Command(arg[0], arg[1:]...).CombinedOutput()
	return strings.TrimSpace(string(t)), err
}

// gather will only download, segment and run the evaluate code, no DB activity
func gather(context *cli.Context) {

	var args []string
	xml := context.Args().First()

	args = []string{context.String("extract"), "SeriesInstanceUID", xml}
	SeriesInstanceUID, err := Run(args...)
	if err != nil || SeriesInstanceUID == "" {
		logger.Fatalf("Could not determine SeriesInstanceUID.  Command is %v", args)
	}

	logger.Info("Processing %v - %v", xml, SeriesInstanceUID)

	// Check the segmented directory
	SegmentedDir := filepath.Join(context.String("segmented"), SeriesInstanceUID)
	ReportFile := filepath.Join(SegmentedDir, "reads.json")
	if Exists(ReportFile) {
		logger.Info("reads file exists, continuing (%v)", ReportFile)
		return
	}
	BaseImage := filepath.Join(SegmentedDir, "image.nii.gz")

	// Download needed?
	DownloadDir := filepath.Join(context.String("dicom"), SeriesInstanceUID)
	if !Exists(DownloadDir) {
		os.MkdirAll(DownloadDir, os.ModePerm|os.ModeDir)
		logger.Info("will be downloading to %v", DownloadDir)
		args = []string{context.String("fetch"), "fetch", "image", "--extract", DownloadDir, SeriesInstanceUID}
		out, err := Run(args...)
		if err != nil {
			logger.Error("error running %v: %v", args, out)
		}
		logger.Info("Fetch: %v %v", args, out)
	} else {
		logger.Info("DICOM exists in %v", DownloadDir)
	}

	if !Exists(SegmentedDir) || !Exists(BaseImage) {
		os.MkdirAll(SegmentedDir, os.ModePerm|os.ModeDir)
		args = []string{context.String("extract"), "segment", xml, DownloadDir, SegmentedDir}
		out, err := Run(args...)
		if err != nil {
			logger.Error("error running %v: %v", args, out)
		}
	} else {
		logger.Info("Segmented exists in %v", SegmentedDir)
	}

	// Copy XML over into the segmented directory
	data, _ := ioutil.ReadFile(xml)
	// Write data to dst
	_, fn := filepath.Split(xml)
	ioutil.WriteFile(filepath.Join(SegmentedDir, fn), data, 0644)

	// Now parse the JSON
	JsonFile := filepath.Join(SegmentedDir, "reads.json")
	fid, err := os.Open(JsonFile)
	if err != nil {
		logger.Fatalf("Error opening %v -- %v", JsonFile, err.Error())
	}
	defer fid.Close()
	everything, _ := jason.NewObjectFromReader(fid)

	// Looping over all the reads
	reads, _ := everything.GetObjectArray("reads")
	for _, read := range reads {
		nodules, _ := read.GetObjectArray("nodules")
		for _, nodule := range nodules {

			read_id, _ := read.GetInt64("id")
			normalized_nodule_id, _ := nodule.GetInt64("normalized_nodule_id")
			centroid, _ := nodule.GetFloat64Array("centroidLPS")
			label_value, _ := nodule.GetFloat64("label_value")
			logger.Debug("Nodule id: %v label: %v centroid: %v", normalized_nodule_id, label_value, centroid)

			// Process the lesion
			GLS := context.String("lesion")

			input := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v.nii.gz", read_id))
			image := filepath.Join(SegmentedDir, "image.nii.gz")

			output := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v.nii.gz", read_id, normalized_nodule_id))
			seed := fmt.Sprintf("%v,%v,%v", centroid[0], centroid[1], centroid[2])
			segmentation_cl := []string{GLS, "-i", image, "-o", output, "--seeds", seed, "--fulloutput"}

			if !Exists(output) {
				out, err := Run(segmentation_cl...)
				logger.Info("running: %v", segmentation_cl)
				if err != nil {
					logger.Error("Error running %v -- %v\nOutput:%v", segmentation_cl, err.Error(), out)
				}
			}

			// evaluate segmentation
			measures := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v_eval.json", read_id, normalized_nodule_id))
			if !Exists(measures) {
				command_line := strings.Trim(fmt.Sprintf("%v", segmentation_cl), "[]")
				args = []string{"python", context.String("evaluate"), "--label", fmt.Sprintf("%v", label_value), "--cli", command_line, output, input, measures}
				logger.Info("Running: %v", args)
				out, err := Run(args...)
				if err != nil {
					logger.Error("Error running %v: %v", args, out)
				}
			}
		}
	}

	// If we are to clean up, delete the DICOM directory
	if context.Bool("clean") || context.Bool("clean-dicom") {
		logger.Debug("Cleaning up %v", DownloadDir)
		os.RemoveAll(DownloadDir)
	}
	if context.Bool("clean") {
		logger.Debug("Removing all images from %v", SegmentedDir)
		fileList, err := filepath.Glob(filepath.Join(SegmentedDir, "*.nii.gz"))
		if err != nil {
			logger.Fatalf("Error globbing %v", err.Error())
		}
		for _, file := range fileList {
			os.Remove(file)
		}
	}
}

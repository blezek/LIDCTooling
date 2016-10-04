package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
)

var ProcessCommand = cli.Command{
	Name:        "process",
	Usage:       "<XML>",
	Description: ProcessDescription,
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
			Name:  "algorithms",
			Value: "algorithms",
			Usage: "Path to directory containing algorithms to run",
		},
		cli.StringFlag{
			Name:  "evaluate",
			Value: "python/evaluateSegmentation.py",
			Usage: "Path to evaluateSegmentation.py program of LIDCTooling",
		},
		cli.StringFlag{
			Name:  "features",
			Value: "python/computeRadiomics.py",
			Usage: "Path to compute radiomics features",
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
	Action: process,
}

func Exists(name string) bool {
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func Run(arg ...string) (string, error) {
	logger.Debug("Starting command " + arg[0])
	cmd := exec.Command(arg[0], arg[1:]...)
	var err error
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b

	// Run in a goroutine, timeout after 10 minutes
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		logger.Debug("Waiting for process")
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(10 * time.Minute):
		logger.Error("Process timeout!")
		err = cmd.Process.Kill()
		err = fmt.Errorf("process timed out")
		break // return b.Bytes(), err
	case err = <-done:
		logger.Debug("Command done " + string(b.Bytes()))
		break // return b.Bytes(), err
	}
	return strings.TrimSpace(string(b.Bytes())), err
}

// process will only download, segment and run the evaluate code, no DB activity
func process(context *cli.Context) {

	var args []string
	xml := context.Args().First()

	args = []string{context.String("extract"), "SeriesInstanceUID", xml}
	SeriesInstanceUID, err := Run(args...)
	if err != nil || SeriesInstanceUID == "" {
		logger.Fatalf("Could not determine SeriesInstanceUID.  Command is %v", args)
	}

	logger.Infof("Processing %v - %v", xml, SeriesInstanceUID)

	// Check the segmented directory
	SegmentedDir := filepath.Join(context.String("segmented"), SeriesInstanceUID)
	ReportFile := filepath.Join(SegmentedDir, "reads.json")
	if Exists(ReportFile) {
		logger.Debugf("reads file exists, continuing (%v)", ReportFile)
		return
	}
	BaseImage := filepath.Join(SegmentedDir, "image.nii.gz")

	// Download needed?
	DownloadDir := filepath.Join(context.String("dicom"), SeriesInstanceUID)
	if !Exists(DownloadDir) {
		os.MkdirAll(DownloadDir, os.ModePerm|os.ModeDir)
		logger.Debugf("will be downloading to %v", DownloadDir)
		args = []string{context.String("fetch"), "fetch", "image", "--extract", DownloadDir, SeriesInstanceUID}
		out, err := Run(args...)
		if err != nil {
			logger.Errorf("error running %v: %v", args, out)
		}
		logger.Debugf("Fetch: %v %v", args, out)
	} else {
		logger.Infof("DICOM exists in %v", DownloadDir)
	}

	if !Exists(SegmentedDir) || !Exists(BaseImage) {
		os.MkdirAll(SegmentedDir, os.ModePerm|os.ModeDir)
		args = []string{context.String("extract"), "segment", xml, DownloadDir, SegmentedDir}
		out, err := Run(args...)
		if err != nil {
			logger.Errorf("error running %v: %v", args, out)
		}
	} else {
		logger.Debugf("Segmented exists in %v", SegmentedDir)
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
			logger.Debugf("Nodule id: %v label: %v centroid: %v", normalized_nodule_id, label_value, centroid)

			// Run the segmentations
			GroundTruth := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v.nii.gz", read_id))
			algorithmDir := context.String("algorithms")
			tag := fmt.Sprintf("_read_%v_nodule_%v", read_id, normalized_nodule_id)
			suffix := tag + ".nii.gz"
			algorithms, _ := ioutil.ReadDir(algorithmDir)
			for _, f := range algorithms {
				if !f.IsDir() {
					algorithm := f.Name()
					outputSegmentation := filepath.Join(SegmentedDir, algorithm+suffix)
					cli := []string{
						filepath.Join(algorithmDir, algorithm),
						"--dicom", DownloadDir,
						"--read", fmt.Sprintf("%v", read_id),
						"--nodule", fmt.Sprintf("%v", normalized_nodule_id),
						"--segmentation_path", SegmentedDir,
						"--ground_truth", GroundTruth,
						"--label_value", fmt.Sprintf("%v", label_value),
						filepath.Join(SegmentedDir, "image.nii.gz"),
						fmt.Sprintf("%v", centroid[0]),
						fmt.Sprintf("%v", centroid[1]),
						fmt.Sprintf("%v", centroid[2]),
						outputSegmentation,
					}
					out, err := Run(cli...)
					logger.Debugf("running: %v", cli)
					logger.Debugf("output: %v", out)
					if err != nil {
						logger.Errorf("Error running %v -- %v\nOutput:%v", cli, err.Error(), out)
					}

					cliString := strings.Join(cli, " ")
					measures := filepath.Join(SegmentedDir, algorithm+tag+".json")
					logger.Debugf("Tag: %v Suffix: %v Measures: %v", tag, suffix, measures)
					args = []string{"python", context.String("evaluate"), "--label", fmt.Sprintf("%v", label_value), "--cli", cliString, outputSegmentation, GroundTruth, measures}
					logger.Debugf("Running evaluation on %v / %v", tag, args)
					out, err = Run(args...)
					if err != nil {
						logger.Errorf("Error running %v: %v", args, out)
					}

					// Calculate PyRadiomics features, assumes the proper Python library is installed
					features := filepath.Join(SegmentedDir, algorithm+"_features"+tag+".json")
					args = []string{"python", context.String("features"),
						"--label", fmt.Sprintf("%v", label_value),
						filepath.Join(SegmentedDir, "image.nii.gz"),
						outputSegmentation, features}
					logger.Infof("Running features on %v / %v", tag, args)
					out, err = Run(args...)
					if err != nil {
						logger.Errorf("Error calculating features\n%v\nError: %v", args, out)
					}

				}
			}

			// 	// Evaluate all the segmentations matching the suffix
			// 	segmentations, _ := filepath.Glob(filepath.Join(SegmentedDir, "*"+suffix))
			// 	for _, match := range segmentations {
			// 		segmentation := filepath.Base(match)
			// 		tag := strings.Split(segmentation, suffix)[0]
			// 		// Process...
			// 		base := basename(basename(suffix))
			// 		measures := filepath.Join(SegmentedDir, tag+base+".json")

			// 		logger.Infof("Tag: %v Suffix: %v Measures: %v", tag, suffix, measures)

			// 		args = []string{"python", context.String("evaluate"), "--label", fmt.Sprintf("%v", label_value), "--cli", tag, match, GroundTruth, measures}
			// 		logger.Infof("Running: %v", args)
			// 		out, err := Run(args...)
			// 		if err != nil {
			// 			logger.Errorf("Error running %v: %v", args, out)
			// 		}
			// 	}
		}
	}

	// If we are to clean up, delete the DICOM directory
	if context.Bool("clean") || context.Bool("clean-dicom") {
		logger.Debugf("Cleaning up %v", DownloadDir)
		os.RemoveAll(DownloadDir)
	}
	if context.Bool("clean") {
		logger.Debugf("Removing all images from %v", SegmentedDir)
		fileList, err := filepath.Glob(filepath.Join(SegmentedDir, "*.nii.gz"))
		if err != nil {
			logger.Fatalf("Error globbing %v", err.Error())
		}
		for _, file := range fileList {
			os.Remove(file)
		}
	}
}

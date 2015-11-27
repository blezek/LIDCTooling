package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
	"github.com/mxk/go-sqlite/sqlite3"
)

var GatherCommand = cli.Command{
	Name:        "gather",
	Usage:       "<XML_1>...<XML_N>",
	Description: GatherDescription,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "db",
			Usage: "SQLite3 database to use, leave blank for none (in-memory only)",
			Value: ":memory:",
		},
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

func gather(context *cli.Context) {

	// Create our tables
	db, err := sqlite3.Open(context.String("db"))
	if err != nil {
		logger.Error("Could not open SQLite3 DB %v: %v", context.String("db"), err.Error())
		os.Exit(1)
	}
	for table, s := range create_tables {
		err = db.Exec(s)
		if err != nil {
			logger.Error("Could not create %v table: %v", table, err.Error())
			os.Exit(1)
		}
	}
	for _, s := range db_migrations {
		err = db.Exec(s)
		if err != nil {
			logger.Error("Could not apply migration (%v): %v", s, err.Error())
			os.Exit(1)
		}
	}

	var args []string
	// Arguments are a list of XML files
	for _, xml := range context.Args() {
		SeriesInstanceUID, _ := Run(context.String("extract"), "SeriesInstanceUID", xml)
		logger.Info("Processing %v - %v", xml, SeriesInstanceUID)

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

		// Check the segmented directory
		SegmentedDir := filepath.Join(context.String("segmented"), SeriesInstanceUID)
		if !Exists(SegmentedDir) {
			os.MkdirAll(SegmentedDir, os.ModePerm|os.ModeDir)
			args = []string{context.String("extract"), "segment", xml, DownloadDir, SegmentedDir}
			out, err := Run(args...)
			if err != nil {
				logger.Error("error running %v: %v", args, out)
			}
		} else {
			logger.Info("Segmented existis in %v", SegmentedDir)
		}

		// Now parse the JSON
		JsonFile := filepath.Join(SegmentedDir, "reads.json")
		fid, err := os.Open(JsonFile)
		if err != nil {
			logger.Error("Error opening %v -- %v", JsonFile, err.Error())
			continue
		}
		defer fid.Close()
		everything, _ := jason.NewObjectFromReader(fid)
		series_uid, _ := everything.GetString("uid")
		Save(db, "series", series_uid, everything)

		// Looping over all the reads
		reads, _ := everything.GetObjectArray("reads")
		for _, read := range reads {
			nodules, _ := read.GetObjectArray("nodules")
			for _, nodule := range nodules {

				// Create an entry for the nodule
				normalized_nodule_id, _ := nodule.GetInt64("normalized_nodule_id")
				nodule_uid := fmt.Sprintf("%v.%v", series_uid, normalized_nodule_id)
				db.Exec("insert or ignore into nodules ( uid, series_uid, normalized_nodule_id ) values (?,?,?)", nodule_uid, series_uid, normalized_nodule_id)

				// Create the read
				read_id, _ := read.GetInt64("id")
				read_uid := fmt.Sprintf("%v.%v.%v", series_uid, read_id, normalized_nodule_id)
				Save(db, "reads", read_uid, nodule)
				characteristics, _ := nodule.GetObject("characteristics")
				Save(db, "reads", read_uid, characteristics)
				db.Exec("update reads set nodule_uid = ? where uid = ?", nodule_uid, read_uid)

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
					args = []string{"python", context.String("evaluate"), "--label", fmt.Sprintf("%v", label_value), output, input, measures}
					logger.Info("Running: %v", args)
					out, err := Run(args...)
					if err != nil {
						logger.Error("Error running %v: %v", args, out)
					}
					// Now, read and combine into a SQLite DB
					fid, _ := os.Open(measures)
					defer fid.Close()
					measure_object, _ := jason.NewObjectFromReader(fid)
					// The UID is the SHA1 of the command string
					command_line := strings.Trim(fmt.Sprintf("%v", segmentation_cl), "[]")
					measures_uid := fmt.Sprintf("%x", sha1.Sum([]byte(command_line)))
					Save(db, "measures", measures_uid, measure_object)
					db.Exec("update measures set command_line = ? where uid = ?", command_line, measures_uid)
					db.Exec("update measures set nodule_uid = ? where uid = ?", nodule_uid, measures_uid)
					db.Exec("update measures set read_uid = ? where uid = ?", read_uid, measures_uid)
				}

			}
		}
	}
}

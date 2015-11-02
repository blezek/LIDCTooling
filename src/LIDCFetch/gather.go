package main

import (
	"crypto/sha1"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
	"github.com/mxk/go-sqlite/sqlite3"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var GatherCommand = cli.Command{
	Name: "gather",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "db",
			Usage: "SQLite3 database to use, leave blank for none (in-memory only)",
			Value: ":memory:",
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

	var args []string
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
		if !Exists(SegmentedDir) {
			os.MkdirAll(SegmentedDir, os.ModePerm|os.ModeDir)
			args = []string{"build/install/LIDCTooling/bin/Extract", "segment", xml, DownloadDir, SegmentedDir}
			Run(args...)
		}

		// Now parse the JSON
		JsonFile := filepath.Join(SegmentedDir, "reads.json")
		fid, _ := os.Open(JsonFile)
		defer fid.Close()
		everything, _ := jason.NewObjectFromReader(fid)
		series_uid, _ := everything.GetString("uid")
		Save(db, "series", series_uid, everything)
		// originalFile, _ := everything.GetString("filename")
		reads, _ := everything.GetObjectArray("reads")
		for _, read := range reads {
			// Looping over all the reads
			read_uid, _ := read.GetString("uid")
			Save(db, "read", read_uid, read)
			db.Exec("update read set series_uid = ? where uid = ?", series_uid, read_uid)
			readId, _ := read.GetInt64("id")
			nodules, _ := read.GetObjectArray("nodules")
			for _, nodule := range nodules {
				nodule_uid, _ := nodule.GetString("uid")
				Save(db, "nodule", nodule_uid, nodule)
				characteristics, _ := nodule.GetObject("characteristics")
				Save(db, "nodule", nodule_uid, characteristics)
				db.Exec("update nodule set read_uid = ? where uid = ?", read_uid, nodule_uid)
				id, _ := nodule.GetInt64("normalized_nodule_id")
				centroid, _ := nodule.GetFloat64Array("centroidLPS")
				label_value, _ := nodule.GetFloat64("label_value")
				logger.Debug("Nodule id: %v label: %v centroid: %v", id, label_value, centroid)

				// Process the lesion
				GLS := "/Users/blezek/Source/ChestImagingPlatform/build/CIP-build/bin/GenerateLesionSegmentation"
				input := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v.nii.gz", readId))
				output := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v.nii.gz", readId, id))
				seed := fmt.Sprintf("%v,%v,%v", centroid[0], centroid[1], centroid[2])
				segmentation_cl := []string{GLS, "-i", input, "-o", output, "--seeds", seed, "--fulloutput"}
				out, err := Run(segmentation_cl...)
				logger.Info("running: %v", segmentation_cl)
				if err != nil {
					logger.Error("Error running %v -- %v\nOutput:%v", segmentation_cl, err.Error(), out)
				}

				// evaluate segmentation
				measures := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v_eval.json", readId, id))
				args = []string{"python", "python/evaluateSegmentation.py", "--label", fmt.Sprintf("%v", label_value), output, input, measures}
				logger.Info("Running: %v", args)
				out, err = Run(args...)
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
				Save(db, "measure", measures_uid, measure_object)
				db.Exec("update measure set command_line = ? where uid = ?", command_line, measures_uid)
				db.Exec("update measure set nodule_uid = ? where uid = ?", nodule_uid, measures_uid)
			}
		}
	}
}

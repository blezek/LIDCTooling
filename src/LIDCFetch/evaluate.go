package main

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
	"github.com/mxk/go-sqlite/sqlite3"
)

var EvaluateCommand = cli.Command{
	Name:        "evaluate",
	Usage:       "<segmented1>...<segmentedN>",
	Description: EvaluateDescription,
	Action:      evaluate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "db",
			Usage: "SQLite3 database to use, leave blank for none (in-memory only)",
			Value: ":memory:",
		},
		cli.StringFlag{
			Name:  "segmented",
			Value: "segmented",
			Usage: "location for segmented files",
		},
	},
}

// The evaluate function collects all the data and stores it in a database
func evaluate(context *cli.Context) {

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
	for idx, SegmentedDir := range context.Args() {

		// Check the segmented directory
		ReportFile := filepath.Join(SegmentedDir, "reads.json")
		if !Exists(ReportFile) {
			logger.Info("reads file DOES NOT exist, skipping (%v)", ReportFile)
			continue
		}

		logger.Debug("Processing (%v/%v) %v", idx+1, len(context.Args()), ReportFile)
		// Now parse the JSON
		JsonFile := filepath.Join(SegmentedDir, "reads.json")
		fid, err := os.Open(JsonFile)
		if err != nil {
			logger.Error("Error opening %v -- %v", JsonFile, err.Error())
			continue
		}
		everything, _ := jason.NewObjectFromReader(fid)
		fid.Close()
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

				// evaluate segmentation
				measures := filepath.Join(SegmentedDir, fmt.Sprintf("read_%v_nodule_%v_eval.json", read_id, normalized_nodule_id))
				if Exists(measures) {

					// Now, read and combine into a SQLite DB
					fid, _ := os.Open(measures)
					measure_object, _ := jason.NewObjectFromReader(fid)
					fid.Close()
					// The UID is the SHA1 of the command string
					command_line, _ := measure_object.GetString("command_line")
					measure_detail, _ := measure_object.GetObject("measures")
					measures_uid := fmt.Sprintf("%x", sha1.Sum([]byte(command_line)))
					Save(db, "measures", measures_uid, measure_detail)
					db.Exec("update measures set command_line = ? where uid = ?", command_line, measures_uid)
					db.Exec("update measures set nodule_uid = ? where uid = ?", nodule_uid, measures_uid)
					db.Exec("update measures set read_uid = ? where uid = ?", read_uid, measures_uid)
				}

			}
		}
	}
}

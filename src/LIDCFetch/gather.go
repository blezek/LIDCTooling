package main

import (
	"encoding/json"
	"github.com/antonholmquist/jason"
	"github.com/codegangsta/cli"
	_ "github.com/jmoiron/jsonq"
	"io/ioutil"
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
	t, err := exec.Command(name, arg...).Output()
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
		fid, err := os.Open(JsonFile)
		var i interface{}
		contents, _ := ioutil.ReadFile(JsonFile)
		json.Unmarshal(contents, &i)
		everything := i.(map[string]interface{})
		reads := everything["reads"].([]interface{})
		for _, r := range reads {
			read := r.(map[string]interface{})
			logger.Debug("reads: %+v", read)
			filename := read["filename"].(string)
			logger.Debug("Filename: %v", filename)
			for _, n := range read["nodules"].([]interface{}) {
				nodule := n.(map[string]interface{})
				id := nodule["id"].(string)
				centroid := make([]float64, 0)
				for _, v := range nodule["centroid"].([]interface{}) {
					centroid = append(centroid, v.(float64))
				}
				label_value := nodule["label_value"].(float64)
				logger.Debug("Nodule id: %v label: %v centroid: %v", id, label_value, centroid)

			}
		}
	}
}

package main

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	cli "github.com/codegangsta/cli"
)

var FetchCommand = cli.Command{
	Name:    "fetch",
	Aliases: []string{},
	Subcommands: []cli.Command{
		FetchImageCommand,
	},
}

var FetchImageCommand = cli.Command{
	Name:        "image",
	Usage:       "<SeriesInstanceUID> [output.zip]",
	Description: "fetch a zip file containing images referenced by SeriesInstanceUID into the output file output.zip if specified, if not specied, save the output in SeriesInstanceUID.zip, in the current directory",
	Action:      fetchImage,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "extract",
			Value: "",
			Usage: "extract the downloaded zip file into the given directory",
		},
	},
}

func fetchImage(context *cli.Context) {
	SeriesInstanceUID := context.Args().First()
	OutputZip := context.Args().Get(1)
	if OutputZip == "" {
		OutputZip = SeriesInstanceUID + ".zip"
	}
	params := url.Values{}
	params.Add("SeriesInstanceUID", SeriesInstanceUID)
	req, err := prepareRequest("/query/getImage", params)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error fetching collections")
	}
	defer resp.Body.Close()

	logger.Debug("Response: %+v", resp)
	if resp.StatusCode != http.StatusOK {
		logger.Error("Got bad status '%v' from server", resp.Status)
		os.Exit(1)
	}

	// Do we extract somewhere?
	if context.String("extract") == "" {
		logger.Debug("Saving to %v", OutputZip)
		f, err := os.Create(OutputZip)
		if err != nil {
			logger.Error("Error opening %v for writing: %v", OutputZip, err.Error())
			return
		}
		defer f.Close()
		io.Copy(f, resp.Body)
	} else {
		outputDir := context.String("extract")
		logger.Debug("Extracting to %v", outputDir)
		err := os.MkdirAll(outputDir, os.ModeDir|os.ModePerm)
		if err != nil {
			logger.Error("Error creating directory %v: %v", outputDir, err.Error())
			return
		}
		fid, err := ioutil.TempFile(os.TempDir(), "LIDCFetch")
		if err != nil {
			logger.Error("Error creating temp file: %v", outputDir, err.Error())
			return
		}
		defer os.Remove(fid.Name())

		// Someday, add a progress bar
		// https://github.com/gosuri/uiprogress
		_, err = io.Copy(fid, resp.Body)
		if err != nil {
			logger.Error("Error downloading zip: %v", err.Error())
			return
		}
		fid.Close()
		logger.Debug("Finished saving to %v", fid.Name())
		// Unzip in the specified location
		z, err := zip.OpenReader(fid.Name())
		if err != nil {
			logger.Error("Error opening zip archive: %v", err.Error())
			return
		}
		defer z.Close()
		// Loop over the files
		for _, f := range z.File {
			extractedFilename := path.Join(outputDir, f.Name)
			logger.Info("Extracting %v Saving %v", f.Name, extractedFilename)
			if strings.HasSuffix(f.Name, "/") {
				// Create the directory
				logger.Debug("Making directory %v", extractedFilename)
				err = os.MkdirAll(extractedFilename, os.ModeDir|os.ModePerm)
				if err != nil {
					logger.Error("Error creating directory %v: %v", extractedFilename, err.Error())
					return
				}
				continue
			}
			efid, err := os.Create(extractedFilename)
			if err != nil {
				logger.Error("Error opening output file %v: %v", extractedFilename, err.Error())
				return
			}
			defer efid.Close()
			zfid, err := f.Open()
			if err != nil {
				logger.Error("Error opening zip file %v: %v", f.Name, err.Error())
				return
			}
			defer zfid.Close()
			_, err = io.Copy(efid, zfid)
			if err != nil {
				logger.Error("Error extracting zip file %v to %v: %v", f.Name, extractedFilename, err.Error())
				return
			}
		}

	}
}

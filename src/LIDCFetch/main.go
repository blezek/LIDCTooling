package main

import (
	"net/http"
	"os"

	cli "github.com/codegangsta/cli"
	log "github.com/op/go-logging"
)

// Valid until 10/2/2016
var APIKey = "25f0025c-071c-426d-b15a-199421e2e889"
var baseURL = "https://services.cancerimagingarchive.net/services/v3/TCIA"
var client = &http.Client{}
var logger = log.MustGetLogger("lidc")

func main() {

	app := cli.NewApp()
	app.Name = "LIDCFetch"
	app.Usage = "Fetch files from LIDC"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose logging",
		},
		cli.StringFlag{
			Name:  "format",
			Value: "json",
			Usage: "format to download",
		},
	}
	app.Commands = []cli.Command{QueryCommand, FetchCommand, GatherCommand}
	app.Before = func(context *cli.Context) error {
		configureLogging(context.Bool("verbose"))
		return nil
	}
	app.Run(os.Args)
}

func configureLogging(verbose bool) {
	backend := log.NewLogBackend(os.Stdout, "", 0)
	f := "%{time:15:04:05.000} %{module} ▶ %{level:.5s} %{id:03x} %{message}"
	f = "%{color}%{time:15:04:05.000} %{module} ▶ %{level:.5s} %{id:03x}%{color:reset} %{message}"
	format := log.MustStringFormatter(f)
	formatter := log.NewBackendFormatter(backend, format)
	log.SetBackend(formatter)
	log.SetLevel(log.CRITICAL, "lidc")
	if verbose {
		log.SetLevel(log.DEBUG, "lidc")
	}
}

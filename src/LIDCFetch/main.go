package main

import (
	"net/http"
	"os"

	cli "github.com/codegangsta/cli"
	log "github.com/op/go-logging"
)

// Valid until 10/2/2016
var APIKey = "864dcc73-ce40-4f19-8a3e-fce71fc2dba2"
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
		cli.BoolFlag{
			Name:  "silent",
			Usage: "no logging, overrides --verbose",
		},
		cli.StringFlag{
			Name:  "format",
			Value: "json",
			Usage: "format to download",
		},
		cli.StringFlag{
			Name:  "apikey",
			Usage: "TCIA API key, to override the internal key",
		},
		cli.StringFlag{
			Name:  "base",
			Usage: "Base URL for TCIA REST API",
			Value: baseURL,
		},
	}
	app.Commands = []cli.Command{QueryCommand, FetchCommand, ProcessCommand, EvaluateCommand}
	app.Before = func(context *cli.Context) error {
		configureLogging(context)
		if context.String("apikey") != "" {
			APIKey = context.String("apikey")
		}
		baseURL = context.String("base")
		return nil
	}
	app.Run(os.Args)
}

func configureLogging(context *cli.Context) {
	backend := log.NewLogBackend(os.Stdout, "", 0)
	f := "%{time:15:04:05.000} %{module} ▶ %{level:.5s} %{id:03x} %{message}"
	//	f = "%{time:15:04:05.000} %{module} ▶ %{level:.5s} %{id:03x} %{message}"
	format := log.MustStringFormatter(f)
	formatter := log.NewBackendFormatter(backend, format)
	log.SetBackend(formatter)
	log.SetLevel(log.INFO, "lidc")
	if context.Bool("verbose") {
		log.SetLevel(log.DEBUG, "lidc")
	}
	if context.Bool("silent") {
		log.SetLevel(log.ERROR, "lidc")
	}
}

package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var GatherCommand = cli.Command{
	Name:   "gather",
	Action: gather,
}

func gather(context *cli.Context) {
	// Walk all the files in LIDC-XML-only
	// var xmlMap map[string]string

	filepath.Walk("LIDC-XML-only", func(path string, info os.FileInfo, err error) error {
		if match, _ := filepath.Match("*.xml", path); match {

		}
		return nil
	})
}

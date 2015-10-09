package main

import cli "github.com/codegangsta/cli"

var QueryCommand = cli.Command{
	Name:    "query",
	Aliases: []string{},
	Subcommands: []cli.Command{
		QueryCollectionCommand,
		//		QuerySeriesCommand,
	},
}

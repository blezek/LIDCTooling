package main

import (
	"io"
	"net/url"
	"os"

	cli "github.com/codegangsta/cli"
)

var QueryCollectionCommand = cli.Command{
	Name:   "collection",
	Usage:  "Query a collection",
	Action: queryCollection,
	Flags:  []cli.Flag{},
}

func queryCollection(context *cli.Context) {
	format := context.GlobalString("format")
	params := url.Values{}
	params.Add("format", format)
	req, err := prepareRequest("/query/getCollectionValues", params)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("Error fetching collections")
	}
	defer resp.Body.Close()
	io.Copy(os.Stdout, resp.Body)
	logger.Debug("Fetching collection %v", resp.Request.URL)
	logger.Debug("Status (%v) %v", resp.StatusCode, resp.Status)
}

package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	if Exists("/this_does_not_exist") {
		t.Error("/this_does_not_exist is claimed to exist")
	}

	wd, _ := os.Getwd()
	if !Exists("/tmp/") {
		t.Error("/tmp/ is claimed to not exist")
	}
	s := filepath.Join(wd, "..", "..", "segmented", "1.3.6.1.4.1.14519.5.2.1.6279.6001.303494235102183795724852353824")
	if !Exists(s) {
		t.Errorf("%v is claimed to not exist", s)
	}
}

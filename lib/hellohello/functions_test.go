package lib

import "testing"

func TestZip(t *testing.T) {
	err := Zip(".testZipDir", "archive.zip")
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

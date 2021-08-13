package lib

import "testing"

func TestZip(t *testing.T) {
	err := Zip("./testZipDir", "archive.zip")
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

func TestUnzip(t *testing.T) {
	idList, err := Unzip("./archive.zip", "hellohello")
	if err != nil {
		t.Errorf("%v\n", err)
	}
	if len(idList) != 1 || idList[0] != "dir" {
		t.Errorf("期待%v，实际%v\n", []string{"dir"}, idList)
	}
}

package main

import (
	"TimeSoft-OA/lib"
	"testing"
)

func TestStoreFile(t *testing.T) {
	err := StoreFile(lib.FileSendHead{
		Name:       ".temp_archive",
		Uploader:   "13284030601",
		ClientCo:   "中石油",
		ScanOrEdit: 0,
	})
	if err != nil {
		t.Errorf("%v\n", err)
	}
}

package filestats

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"
)

func TestNew(t *testing.T) {
	got := New()
	want := &Stats{}
	if reflect.TypeOf(got) != reflect.TypeOf(want) {
		t.Errorf("Expected %+v got %+v", want, got)
	}
}

func TestAddFile(t *testing.T) {
	stats := New()

	testZip := filepath.Join("../../", "testdata", "zips", "filestats.zip")
	archive, err := ioutil.ReadFile(testZip)
	if err != nil {
		t.Errorf("Could not open Test file: %s\n", testZip)
	}

	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		t.Errorf("Could not decode Test file: %s\n", testZip)
	}

	files := [][]string{
		{"test.php", ".php", "87"},
		{"test.css", ".css", "58"},
	}

	for _, f := range zr.File {
		stats.AddFile(f)
	}

	for k, want := range files {
		got := stats.Files[k]
		size, _ := strconv.Atoi(want[2])
		if got.Name != want[0] || got.Extension != want[1] || got.Size != int64(size) {
			t.Errorf("Expected %+v got %+v", want, got)
		}
	}
}

func TestGenerateSummary(t *testing.T) {
	stats := New()

	testZip := filepath.Join("../../", "testdata", "zips", "filestats.zip")
	archive, err := ioutil.ReadFile(testZip)
	if err != nil {
		t.Errorf("Could not open Test file: %s\n", testZip)
	}

	zr, err := zip.NewReader(bytes.NewReader(archive), int64(len(archive)))
	if err != nil {
		t.Errorf("Could not decode Test file: %s\n", testZip)
	}

	for _, f := range zr.File {
		stats.AddFile(f)
	}

	stats.GenerateSummary()

	wantFiles := 2
	gotFiles := stats.TotalFiles
	if wantFiles != gotFiles {
		t.Errorf("Expected %d got %d", wantFiles, gotFiles)
	}

	wantSize := int64(145)
	gotSize := stats.TotalSize
	if wantSize != gotSize {
		t.Errorf("Expected %d got %d", wantSize, gotSize)
	}

	wantPHP := uint8(60)
	gotPHP := stats.Summary.PHP
	if wantPHP != gotPHP {
		t.Errorf("Expected %d got %d", wantPHP, gotPHP)
	}

	wantJS := uint8(0)
	gotJS := stats.Summary.JS
	if wantJS != gotJS {
		t.Errorf("Expected %d got %d", wantJS, gotJS)
	}

	wantCSS := uint8(40)
	gotCSS := stats.Summary.CSS
	if wantCSS != gotCSS {
		t.Errorf("Expected %d got %d", wantCSS, gotCSS)
	}
}

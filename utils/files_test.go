package utils

import (
	"path/filepath"
	"testing"
)

func TestIsIgnored(t *testing.T) {
	names := []string{
		"node_modules",
		".git",
		"vendor",
		"test.db",
		"yarn.lock",
	}

	for _, n := range names {
		if !IsIgnored(n) {
			t.Errorf("File/folder %s wasn't ignored, but should have been", n)
		}
	}
}

func TestIsDir(t *testing.T) {
	exists, err := filepath.Abs("../commands")
	if err != nil {
		t.Fatal(err)
	}

	if !IsDir(exists) {
		t.Errorf("Returned false on existing directory: %s", exists)
	}

	doesNotExists, err := filepath.Abs("../fake")
	if err != nil {
		t.Fatal(err)
	}

	if IsDir(doesNotExists) {
		t.Errorf("Returned true on directory that doesn't exist: %s", doesNotExists)
	}
}

func TestIsFile(t *testing.T) {
	exists, err := filepath.Abs("./files.go")
	if err != nil {
		t.Fatal(err)
	}

	if !IsFile(exists) {
		t.Errorf("Returned false on existing file: %s", exists)
	}

	doesNotExists, err := filepath.Abs("../fake.go")
	if err != nil {
		t.Fatal(err)
	}

	if IsFile(doesNotExists) {
		t.Errorf("Returned true on file that doesn't exist: %s", doesNotExists)
	}
}

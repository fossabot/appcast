package file_test

import (
	"path/filepath"
	"testing"

	"github.com/abemedia/appcast/internal/test"
	"github.com/abemedia/appcast/target/file"
)

func TestFile(t *testing.T) {
	path := t.TempDir()

	tgt, err := file.New(file.Config{Path: path})
	if err != nil {
		t.Fatal(err)
	}

	test.Target(t, tgt, func(asset string) string {
		return "file://" + filepath.Join(path, asset)
	})
}

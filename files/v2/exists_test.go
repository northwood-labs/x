package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileExistsFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if !FileExists(path) {
		t.Fatalf("expected file to exist")
	}
}

func TestFileExistsMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.txt")

	if FileExists(path) {
		t.Fatalf("expected missing file to return false")
	}
}

func TestFileExistsDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	if FileExists(dir) {
		t.Fatalf("expected directory path to return false")
	}
}

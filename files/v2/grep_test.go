package files

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGrepFileFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("hello world\nthis is a test\n"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	found, err := GrepFile(path, "hello")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !found {
		t.Fatalf("expected match, got no match")
	}
}

func TestGrepFileNotFound(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("hello world\nthis is a test\n"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	found, err := GrepFile(path, "nomatch")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found {
		t.Fatalf("expected no match, got match")
	}
}

func TestGrepFileInvalidRegex(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("hello world\n"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	found, err := GrepFile(path, "[")
	if err == nil {
		t.Fatalf("expected regex compile error, got nil")
	}
	if found {
		t.Fatalf("expected found=false on error, got true")
	}
}

func TestGrepFileMissingFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.txt")

	found, err := GrepFile(path, "hello")
	if err == nil {
		t.Fatalf("expected read error, got nil")
	}
	if found {
		t.Fatalf("expected found=false on error, got true")
	}
}

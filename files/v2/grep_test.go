// Copyright 2026, Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		t.Fatal("expected match, got no match")
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
		t.Fatal("expected no match, got match")
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
		t.Fatal("expected regex compile error, got nil")
	}

	if found {
		t.Fatal("expected found=false on error, got true")
	}
}

func TestGrepFileMissingFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.txt")

	found, err := GrepFile(path, "hello")
	if err == nil {
		t.Fatal("expected read error, got nil")
	}

	if found {
		t.Fatal("expected found=false on error, got true")
	}
}

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

// TestFileExistsFile confirms the happy path: a real file on disk is
// reported as existing.
func TestFileExistsFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sample.txt")

	if err := os.WriteFile(path, []byte("data"), 0o600); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	if !FileExists(path) {
		t.Fatal("expected file to exist")
	}
}

// TestFileExistsMissing verifies that a non-existent path returns
// false rather than panicking or returning an error.
func TestFileExistsMissing(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "does-not-exist.txt")

	if FileExists(path) {
		t.Fatal("expected missing file to return false")
	}
}

// TestFileExistsDirectory ensures that a directory is not mistaken for
// a file. Callers rely on this distinction to avoid accidentally
// reading a directory as if it were a config file.
func TestFileExistsDirectory(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	if FileExists(dir) {
		t.Fatal("expected directory path to return false")
	}
}

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

import "os"

// FileExists reports whether a regular file (not a directory) exists at
// the given path. CLI tools call this before reading or overwriting a
// file to provide clear error messages ("file not found") rather than
// letting downstream os.Open failures bubble up with less context.
// Directories are excluded because callers that care about a directory
// use os.Stat directly — conflating the two would mask bugs where a
// path accidentally points to a directory.
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

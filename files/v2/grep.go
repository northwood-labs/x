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
	"fmt"
	"os"
	"regexp"
)

// GrepFile searches the contents of a file for a regex pattern and
// returns whether a match was found. This is used by config validation
// to confirm that managed markers or expected content are present in a
// file without requiring the caller to handle file I/O and regex
// compilation separately. The entire file is read into memory because
// config files are small and line-by-line scanning would complicate
// multi-line pattern matching.
func GrepFile(path, s string) (bool, error) {
	b, err := os.ReadFile(path) // lint:allow_possible_insecure
	if err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	found, err := regexp.Match(s, b)
	if err != nil {
		return false, fmt.Errorf("could not match regexp pattern: %w", err)
	}

	return found, nil
}

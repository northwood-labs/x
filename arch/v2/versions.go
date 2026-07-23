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

package arch

import "fmt"

// wasmID is extracted as a constant because it appears in multiple
// conditional checks below and a typo would silently produce the wrong
// platform string.
const wasmID = "wasm"

// GetFriendlyName takes the values of GOOS and GOARCH, and returns a
// human-readable platform name (U.S. English). This is the single entry
// point for turning raw build-target identifiers into something suitable
// for version banners, release notes, and user-facing diagnostics.
func GetFriendlyName(osStr, archStr string) string {
	var (
		osFriendly   string
		archFriendly string
	)

	osFriendly = osStr
	archFriendly = archStr

	// WebAssembly targets are special-cased because combining OS + arch
	// ("JavaScript on WebAssembly") reads awkwardly. A single branded
	// name is clearer for users.
	if osStr == "js" && archStr == wasmID {
		return "WebAssembly"
	}

	if osStr == "wasip1" && archStr == wasmID {
		return "WebAssembly with WASI Preview 1"
	}

	// Apple Silicon is the marketing name users recognize; "macOS on
	// ARM (64-bit)" would be technically correct but unfamiliar.
	if osStr == "darwin" && archStr == "arm64" {
		return "macOS on Apple Silicon"
	}

	// Fall through to the generic lookup tables for all other
	// combinations, keeping the raw identifier as a fallback if the
	// map doesn't contain the value (future Go ports).
	if val, ok := OSMap[osStr]; ok {
		osFriendly = val
	}

	if val, ok := ArchMap[archStr]; ok {
		archFriendly = val
	}

	return fmt.Sprintf("%s on %s", osFriendly, archFriendly)
}

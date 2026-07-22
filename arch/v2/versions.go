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

const wasmID = "wasm"

// GetFriendlyName takes the values of GOOS and GOARCH, and returns a "friendly"
// name for the pairing (U.S. English).
func GetFriendlyName(osStr, archStr string) string {
	var (
		osFriendly   string
		archFriendly string
	)

	osFriendly = osStr
	archFriendly = archStr

	if osStr == "js" && archStr == wasmID {
		return "WebAssembly"
	}

	if osStr == "wasip1" && archStr == wasmID {
		return "WebAssembly with WASI Preview 1"
	}

	if osStr == "darwin" && archStr == "arm64" {
		return "macOS on Apple Silicon"
	}

	if val, ok := OSMap[osStr]; ok {
		osFriendly = val
	}

	if val, ok := ArchMap[archStr]; ok {
		archFriendly = val
	}

	return fmt.Sprintf("%s on %s", osFriendly, archFriendly)
}

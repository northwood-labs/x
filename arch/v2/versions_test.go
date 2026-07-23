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

import "testing"

// errorMessage is reused across all table-driven assertions to provide
// consistent failure output showing expected vs actual values.
const errorMessage = "Result was `%s` instead of `%s`."

// TestGetFriendlyNameJS confirms the pure-WebAssembly special case
// (js/wasm) collapses into a single branded name rather than the
// generic "OS on Arch" format.
func TestGetFriendlyNameJS(t *testing.T) {
	expected := "WebAssembly"
	actual := GetFriendlyName("js", "wasm")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

// TestGetFriendlyNameASi confirms the Apple Silicon special case
// produces the marketing name users expect rather than "macOS on
// ARM (64-bit)".
func TestGetFriendlyNameASi(t *testing.T) {
	expected := "macOS on Apple Silicon"
	actual := GetFriendlyName("darwin", "arm64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

// TestGetFriendlyNameLinux64 exercises the common-path lookup where
// both OS and arch are present in their respective maps.
func TestGetFriendlyNameLinux64(t *testing.T) {
	expected := "Linux on Intel (64-bit)"
	actual := GetFriendlyName("linux", "amd64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

// TestGetFriendlyNameNoOS verifies that a less common OS (illumos)
// still resolves correctly when present in OSMap.
func TestGetFriendlyNameNoOS(t *testing.T) {
	expected := "illumos on Intel (64-bit)"
	actual := GetFriendlyName("illumos", "amd64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

// TestGetFriendlyNameNoOSArch verifies the mainframe combination
// (zos/s390x) maps to full branded names for both OS and arch.
func TestGetFriendlyNameNoOSArch(t *testing.T) {
	expected := "IBM z/OS on System/390 (64-bit)"
	actual := GetFriendlyName("zos", "s390x")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

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

// Bit-width and vendor-family constants are composed into the ArchMap
// values below. Keeping them as named constants avoids repetition and
// makes it easy to update branding in one place (e.g., if "Loongson"
// were rebranded).
const (
	b32 = " (32-bit)"
	b64 = " (64-bit)"

	intel = "Intel"
	arm   = "ARM"
	loong = "Loongson"
	mips  = "MIPS"
	ppc   = "PowerPC"
	riscv = "RISC-V"
	s390  = "System/390"
	sparc = "SPARC"
	wasm  = "WebAssembly"
)

// go tool dist list.

var (
	// OSMap translates GOOS identifiers into user-facing operating system
	// names. The Go toolchain uses terse identifiers that are meaningful to
	// developers but opaque to end-users reading release notes, version
	// strings, or diagnostic output.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// http://bit.ly/4b8bHz3
	OSMap = map[string]string{
		"aix":       "IBM AIX",        // https://www.ibm.com/products/aix
		"android":   "Android",        // https://www.android.com
		"darwin":    "macOS",          // https://www.apple.com/macos
		"dragonfly": "DragonFly BSD",  // https://www.dragonflybsd.org
		"freebsd":   "FreeBSD",        // https://www.freebsd.org
		"hurd":      "GNU Hurd",       // https://www.gnu.org/software/hurd, Not enabled by default
		"illumos":   "illumos",        // https://www.illumos.org
		"ios":       "iOS",            // https://www.apple.com/ios
		"js":        "JavaScript",     // https://ecma-international.org/publications-and-standards/standards/ecma-262/
		"linux":     "Linux",          // https://docs.kernel.org
		"netbsd":    "NetBSD",         // https://www.netbsd.org
		"openbsd":   "OpenBSD",        // https://www.openbsd.org
		"plan9":     "Plan 9 UNIX",    // https://9p.io/plan9
		"solaris":   "Solaris",        // https://www.oracle.com/solaris
		"wasip1":    "WASI Preview 1", // https://go.dev/blog/wasi
		"windows":   "Windows",        // https://windows.com
		"zos":       "IBM z/OS",       // https://www.ibm.com/products/zos
	}

	// ArchMap translates GOARCH identifiers into friendly architecture
	// names that combine the vendor family with the bit-width. This lets
	// downstream code produce strings like "ARM (64-bit)" without
	// hard-coding display logic at each call site.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// http://bit.ly/4b8bHz3
	ArchMap = map[string]string{
		"386":         intel + b32,
		"amd64":       intel + b64,
		"amd64p32":    intel + b64,
		"arm":         arm + b32,
		"arm64":       arm + b64,
		"arm64be":     arm + b64,
		"armbe":       arm + b32,
		"loong64":     loong + b64,
		"mips":        mips + b32,
		"mips64":      mips + b64,
		"mips64le":    mips + b64,
		"mips64p32":   mips + b64,
		"mips64p32le": mips + b64,
		"mipsle":      mips + b32,
		"ppc":         ppc + b32,
		"ppc64":       ppc + b64,
		"ppc64le":     ppc + b64,
		"riscv":       riscv + b32,
		"riscv64":     riscv + b64,
		"s390":        s390 + b32,
		"s390x":       s390 + b64,
		"sparc":       sparc + b32,
		"sparc64":     sparc + b64,
		"wasm":        wasm,
	}
)

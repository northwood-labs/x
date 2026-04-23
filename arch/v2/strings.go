package arch

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

// go tool dist list

var (
	// OSMap is a mapping of GOOS values to "friendly" values.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// https://github.com/golang/go/blob/master/src/go/build/syslist.go
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

	// ArchMap is a mapping of GOARCH values to "friendly" values.
	// https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
	// https://github.com/golang/go/blob/master/src/go/build/syslist.go
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

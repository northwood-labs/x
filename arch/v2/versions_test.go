package arch

import "testing"

const errorMessage = "Result was `%s` instead of `%s`."

func TestGetFriendlyNameJS(t *testing.T) {
	expected := "WebAssembly"
	actual := GetFriendlyName("js", "wasm")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

func TestGetFriendlyNameASi(t *testing.T) {
	expected := "macOS on Apple Silicon"
	actual := GetFriendlyName("darwin", "arm64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

func TestGetFriendlyNameLinux64(t *testing.T) {
	expected := "Linux on Intel (64-bit)"
	actual := GetFriendlyName("linux", "amd64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

func TestGetFriendlyNameNoOS(t *testing.T) {
	expected := "illumos on Intel (64-bit)"
	actual := GetFriendlyName("illumos", "amd64")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

func TestGetFriendlyNameNoOSArch(t *testing.T) {
	expected := "IBM z/OS on System/390 (64-bit)"
	actual := GetFriendlyName("zos", "s390x")

	if actual != expected {
		t.Errorf(errorMessage, actual, expected)
	}
}

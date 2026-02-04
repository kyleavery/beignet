package stager

import (
	"errors"
	"runtime"
)

var ErrUnsupportedPlatform = errors.New("stager: unsupported platform/arch")

// LoaderText returns the embedded darwin/arm64 loader image bytes and the entry
// function offset (relative to the start of the image).
func LoaderText() ([]byte, uint64, error) {
	if runtime.GOOS == "darwin" && runtime.GOARCH == "arm64" {
		text, entryOff := loaderTextDarwinArm64()
		return text, entryOff, nil
	}
	return nil, 0, ErrUnsupportedPlatform
}

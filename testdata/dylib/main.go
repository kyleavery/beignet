package main

/*
#include <stdint.h>
*/
import "C"

import (
	"os"
)

const markerPath = "/tmp/beignet_test_marker"

//export BeignetEntry
func BeignetEntry() {
	_ = os.WriteFile(markerPath, []byte("ok"), 0o600)
}

func main() {}

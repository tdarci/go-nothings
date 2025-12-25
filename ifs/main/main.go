package main

import (
	"runtime"

	"github.com/tdarci/go-nothings/ifs/harness"
)

func main() {
	runtime.LockOSThread()
	harness.RunSierpinski()
}
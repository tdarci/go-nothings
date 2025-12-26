package main

import (
	"log"
	"runtime"
	"time"

	"github.com/tdarci/go-nothings/ifs/harness"
)

const (
	maxIterations = 100_000
	drawingDelay = time.Millisecond * 2
	containerEdgeLen = 1_000
	maxFernW = 100
	maxFernH = 100
)

func main() {

	st := time.Now()
	timeStr := st.Format("2006-01-02_15-04-05")
	// filename := "ifs.log"
	// filename := timeStr + "-ifs.log"
	// f, err := os.Create(filename)
	// if err != nil {
	// 	log.Fatalf("unable to create file %q: %s", filename, err)
	// }
	// log.SetOutput(f)
	log.Printf("Starting IFS processing at %s...", timeStr)

	runtime.LockOSThread()
	sys := harness.NewSierpinskiSystem(containerEdgeLen)
	// sys := harness.NewFernSystem(maxFernW, maxFernH)
	harness.Run(sys, maxIterations, drawingDelay)
}

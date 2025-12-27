package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/tdarci/go-nothings/ifs/harness"
	"github.com/tdarci/go-nothings/ifs/shared"
)

const (
	maxIterations    = 100_000
	drawingDelay     = time.Millisecond * 2
	containerEdgeLen = 1_000
	maxFernW         = 1000
	maxFernH         = 1000
)

var (
	minPt = shared.Point{X: -100.0, Y: -150.0}
	maxPt = shared.Point{X: 300.0, Y: 400.0}
)

func main() {

	usage := func() {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("=====================================")
		fmt.Println("Usage:")
		fmt.Println(" ifs triangle|fern|rectangle")
	}

	if len(os.Args) < 3 {
		usage()
		return
	}

	st := time.Now()
	timeStr := st.Format("2006-01-02_15-04-05")
	// filename := "ifs.log"
	// filename := timeStr + "-ifs.log"
	// f, err := os.Create(filename)
	// if err != nil {
	// 	log.Fatalf("unable to create file %q: %s", filename, err)
	// }
	// log.SetOutput(f)

	var sys *harness.System[shared.Point]

	fractal := os.Args[2]
	switch fractal {
	case "triangle":
		sys = harness.NewSierpinskiSystem(containerEdgeLen)
	case "fern":
		sys = harness.NewFernSystem(maxFernW, maxFernH)
	case "rectangle":
		sys = harness.NewNegTestSystem(minPt, maxPt)
	default:
		usage()
		return
	}

	log.Printf("Starting IFS processing of %s at %s...", fractal, timeStr)

	runtime.LockOSThread()

	harness.Run(sys, maxIterations, drawingDelay)
}

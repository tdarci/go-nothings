package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/tdarci/go-nothings/ifs/harness"
)

const (
	iterations = 10_000
)

func main() {

	st := time.Now()
	timeStr := st.Format("2006-01-02_15-04-05")
	filename := "ifs.log"
	// filename := timeStr + "-ifs.log"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("unable to create file %q: %s", filename, err)
	}
	log.SetOutput(f)
	log.Printf("Starting log for %s...", timeStr)

	runtime.LockOSThread()
	harness.RunSierpinski(iterations)
}

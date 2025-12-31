package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/tdarci/go-nothings/ifs/config"
	"github.com/tdarci/go-nothings/ifs/harness"
	"github.com/tdarci/go-nothings/ifs/shared"
	"gopkg.in/yaml.v3"
)

var (
	defaultConfig = config.Config{
		MaxIterations:  100_000,
		DrawingDelayMs: 2,
		Triangle: config.TriangleConfig{
			EdgeLen: 1_000.0,
		},
		Carpet: config.CarpetConfig{
			EdgeLen: 1_000.0,
		},
		Fern: config.FernConfig{
			MaxWidth:  1000,
			MaxHeight: 1000,
			Thresholds: config.FernThresholds{
				StemBaseThresholdPct:       1,
				SmallerLeafletThresholdPct: 86,
				LeftBigLeafletThresholdPct: 93,
			},
		},
		Rectangle: config.RectangleConfig{
			MinPoint: shared.Point{X: -200, Y: -300},
			MaxPoint: shared.Point{X: 100, Y: 200},
		},
	}
)

func main() {

	usage := func() {
		fmt.Println("")
		fmt.Println("")
		fmt.Println("=====================================")
		fmt.Println("Usage:")
		fmt.Println(" ifs triangle|carpet|fern|rectangle [configFile]")
	}

	if len(os.Args) < 3 {
		usage()
		return
	}

	st := time.Now()

	var configFileName string
	var cfg config.Config = defaultConfig
	if len(os.Args) > 3 {
		configFileName = os.Args[3]
		data, err := os.ReadFile(configFileName)
		if err != nil {
			fmt.Printf("Unable to read file %q: %s", configFileName, err)
			return
		}
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			fmt.Printf("Error reading yaml file %q into config: %s", configFileName, err)
			return
		}
	}

	var sys *harness.System[shared.Point]

	fractal := os.Args[2]
	switch fractal {
	case "triangle":
		sys = harness.NewSierpinskiSystem(cfg.Triangle)
	case "carpet":
		sys = harness.NewCarpetSystem(cfg.Carpet)
	case "fern":
		sys = harness.NewFernSystem(cfg.Fern)
	case "rectangle":
		sys = harness.NewNegTestSystem(cfg.Rectangle)
	default:
		usage()
		return
	}

	timeStr := st.Format("2006-01-02_15-04-05")
	logFilename := "ifs.log"
	log.Printf("==========\nStarting iterated function system processing of %q at %s. Logs are being written to %q\n\n", fractal, timeStr, logFilename)

	f, err := os.Create(logFilename)
	if err != nil {
		log.Fatalf("unable to create file %q: %s", logFilename, err)
	}
	log.SetOutput(f)

	log.Printf("Starting iterated function system processing of %q at %s...", fractal, timeStr)

	runtime.LockOSThread()

	harness.Run(sys, cfg.MaxIterations, time.Millisecond*time.Duration(cfg.DrawingDelayMs))
}

package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/tdarci/go-nothings/ifs/harness"
	"github.com/tdarci/go-nothings/ifs/shared"
	"gopkg.in/yaml.v3"
)

var (
	defaultConfig = Config{
		MaxIterations:  100_000,
		DrawingDelayMs: 2,
		Triangle: TriangleConfig{
			EdgeLen: 1_000.0,
		},
		Fern: FernConfig{
			MaxWidth:  1000,
			MaxHeight: 1000,
		},
		Rectangle: RectangleConfig{
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

	var configFileName string
	var cfg Config = defaultConfig
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

	log.Printf("FOOBAR config: %+v", cfg)

	var sys *harness.System[shared.Point]

	fractal := os.Args[2]
	switch fractal {
	case "triangle":
		sys = harness.NewSierpinskiSystem(cfg.Triangle.EdgeLen)
	case "fern":
		sys = harness.NewFernSystem(cfg.Fern.MaxWidth, cfg.Fern.MaxHeight)
	case "rectangle":
		sys = harness.NewNegTestSystem(cfg.Rectangle.MinPoint, cfg.Rectangle.MaxPoint)
	default:
		usage()
		return
	}

	log.Printf("Starting IFS processing of %s at %s...", fractal, timeStr)

	runtime.LockOSThread()

	harness.Run(sys, cfg.MaxIterations, time.Millisecond*time.Duration(cfg.DrawingDelayMs))
}

type Config struct {
	MaxIterations  int             `yaml:"MaxIterations"`
	DrawingDelayMs int             `yaml:"DrawingDelayMs"`
	Triangle       TriangleConfig  `yaml:"Triangle"`
	Fern           FernConfig      `yaml:"Fern"`
	Rectangle      RectangleConfig `yaml:"Rectangle"`
}

type TriangleConfig struct {
	EdgeLen float64 `yaml:"EdgeLen"`
}

type FernConfig struct {
	MaxWidth  float64 `yaml:"MaxWidth"`
	MaxHeight float64 `yaml:"MaxHeight"`
}

type RectangleConfig struct {
	MinPoint shared.Point `yaml:"MinPoint"`
	MaxPoint shared.Point `yaml:"MaxPoint"`
}

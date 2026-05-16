package main

import (
	"fmt"
	"os"

	"github.com/tdarci/go-nothings/quakes-roster/internal/pdf"
	"github.com/tdarci/go-nothings/quakes-roster/internal/scraper"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  scrape   - fetch players + images")
		fmt.Println("  pdf      - generate PDF")
		return
	}

	switch os.Args[1] {
	case "scrape":
		scraper.Run()
	case "pdf":
		pdf.Run()
	default:
		fmt.Println("Unknown command")
	}
}

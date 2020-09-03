package main

// Thanks for the assist, Kinya El Grande

// TRY IT OUT
// ====================================================================================================================================================
//     go run main.go --infile="https://www.eastbaytimes.com/wp-content/uploads/2016/08/20070511_123707_tower1.jpg?w=476" > /tmp/building.jpg
//     go run main.go --infile="https://images.wagwalkingweb.com/media/articles/dog/why-is-my-dog-jumping/why-is-my-dog-jumping.jpg" > /tmp/dog.jpg
// ====================================================================================================================================================

import (
	"bufio"
	"flag"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"

	"./filters"
)

func main() {
	infile := flag.String("infile", "", "please supply a url")
	flag.Parse()
	var in io.Reader
	if infile != nil && *infile != "" {
		x, downloadErr := DownloadURL(*infile)
		if downloadErr != nil {
			log.Panicf("download blew up: %s", downloadErr)
		}
		defer x.Close()
		in = x
	} else {
		in = bufio.NewReader(os.Stdin)
	}

	//decode  Image
	imgContent, imgFormat, err := image.Decode(in)
	if err != nil {
		panic(err.Error())
	}

	filtered := filters.Redden(imgContent)

	if imgFormat == "jpeg" {
		jpeg.Encode(os.Stdout, filtered, nil)
	} else if imgFormat == "png" {
		png.Encode(os.Stdout, filtered)
	} else {
		log.Panicln("File format must be png or jpg")
	}
}

//DownloadURL downloads from the specified url
func DownloadURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const photoCount = 22

func main() {
	c := spewOutFiles()
	size := generateThumbnails(c)
	fmt.Printf("THIS IS THE SIZE ----> %d <------\n", size)

}

type funFile struct {
	Name string
	Size int64
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// spewOutFiles is a helper function that returns a channel immediately.
// this channel is populated over time with photoCount files
func spewOutFiles() <-chan funFile {
	c := make(chan funFile)
	go func() {
		for i := 0; i < photoCount; i++ {
			time.Sleep(time.Millisecond * time.Duration(rand.Int63n(100)+1))
			ff := funFile{
				Name: fmt.Sprintf("photo_%02d.jpg", i),
				Size: rand.Int63n(1000) + 1,
			}
			log.Printf("%s: SPEWED", ff.Name)
			c <- ff
		}
		close(c)
	}()
	return c
}

// makeAThumbnail is function that acts as a stand-in for the process of creating a thumbnail of a file
func makeAThumbnail(f funFile) (funFile, error) {
	time.Sleep(time.Millisecond * time.Duration(rand.Int63n(1000)+1))
	f.Size = f.Size / 2
	return f, nil
}

// generateThumbnails takes an inbound channel of files, generates a thumbnail for each one, and returns the cumulative size
// of all the files once the channel is closed
func generateThumbnails(filenames <-chan funFile) int64 {
	sizes := make(chan int64)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for f := range filenames {
			wg.Add(1)
			// worker
			go func(infile funFile) {
				defer wg.Done()
				log.Printf("%s: PROCESSING", infile.Name)
				thumb, err := makeAThumbnail(infile)
				log.Printf("%s: PROCESSED", thumb.Name)
				if err != nil {
					log.Println(err)
					return
				}
				sizes <- thumb.Size
			}(f)
		}
		wg.Done() // filenames is now closed
	}()

	// closer
	go func() {
		wg.Wait()
		close(sizes)
	}()

	var total int64
	for size := range sizes {
		total += size // total is protected since it is only written by the path through this unbuffered channel
	}
	return total
}

// implements web crawler from chapter 8 of The Go Programming Language, with termination
// This code is written to run on https://play.golang.org/
//
// In order to change the output, change the value of randomSeed.
// You can also experiment with concurrentWorkerCount.
//
package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	randomSeed            = 123455845 // change this value to have the program traverse in a different order when running on https://play.golang.org/
	concurrentWorkerCount = 10
)

var (
	startUrls = []string{"http://start-foo.bar/wookie"}
	urls      = []string{
		"http://a-aardvarks.com/momo",
		"http://b-recipes.com/gumbo",
		"http://c-clowns.com/bozo",
		"http://d-dogs.com/rambo",
		"http://e-errors.org/ohno",
		"http://f-fauxpas.com/nono",
		"http://g-garbage.com/packaging",
		"http://h-hotwheels.com/speed-racer",
		"http://i-incredibles.com/edna-mode",
	}
)

func init() {
	rand.Seed(randomSeed + time.Now().UnixNano())
}

// programArgs() returns the arguments supplied to our program.
// In this case we are supplying a hard-coded value.
func programArgs() []string {
	// return os.Args[1:] <-- this is what we'd really use
	return startUrls
}

// extract() is a test method that returns a set of random urls.
// in the real world, extract() would scan the provided url for links.
func extract(url string) (foundUrls []string) {

	lowVal := rand.Int31n(int32(len(urls)))
	interval := int32(len(urls)-1) - lowVal
	// log.Printf("low: %d. interval: %d", lowVal, interval)
	highVal := lowVal + 1
	if interval > 0 {
		highVal += rand.Int31n(interval)
	}

	foundUrls = urls[lowVal:highVal]
	return
}

func main() {
	// setup channels
	crawlResults := make(chan []string) // channel of lists of urls to process. may have duplicates
	unseenLinks := make(chan string)    // channel of de-duped urls
	pending := make(chan int)           // channel of things waiting to be completed
	pendCount := 0

	// process input
	go func() {
		pending <- 1                  // adding item to crawlResults
		crawlResults <- programArgs() // this gets things rolling...
	}()

	// Notice when we have no more pending operations.
	// At that point, this program is done.
	go func() {
		// Once this function is running it will take items off of pending as soon as they are put there. So no blocking to worry about.
		for pend := range pending {
			pendCount += pend
			// log.Printf("Pending: %d", pendCount)
			if pendCount <= 0 {
				// all done!
				close(crawlResults) // ** This will end our program. **
			}
		}
	}()

	// spawn crawlers (workers that run extract())
	for i := 0; i < concurrentWorkerCount; i++ {
		// each of these makes a crawler to run concurrently. these goroutines die when the program exits
		go func() {
			for link := range unseenLinks { // We will block here until there's something to take off of the channel and keep looping until something closes unseenLinks or main exits
				//
				// We have 2 actions for pending channel now... +1 for initiating extract() and -1 for taking something off of unseenLinks. So no need to do anything.
				foundLinks := extract(link)
				if foundLinks != nil {
					pending <- 1 // adding item to crawlResults
					go func() { crawlResults <- foundLinks }()
				}
				pending <- -1 // extract() is complete
			}
		}()
	}

	// process items in crawlResults. put new urls on unseenLinks
	seen := make(map[string]bool)
	for list := range crawlResults { // when crawlResults is closed, this loop is done and our program finishes (and closes any open goroutines)
		for _, curLink := range list {
			if !seen[curLink] {
				// Print out every time we find one.
				fmt.Printf("* Found a link: %s\n", curLink)
				seen[curLink] = true
				pending <- 1 // adding item to unseenLinks
				unseenLinks <- curLink
			}
		}
		pending <- -1 // item taken off of crawlResults
	}
}

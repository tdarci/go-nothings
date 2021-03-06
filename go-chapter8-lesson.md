### Chapter 8: Goroutines & Channels

Computers are amazing... they can do more than one thing at a time... this laptop has several CPUs... and each can be doing a different task.

So, instead of saying:

1. do this
2. then do that
3. and finally do the third thing

We can say this:
- do this
- do that
- do the third thing

and let the computer do them _all at once_

however, this can get tricky when we want to _coordinate_ between these different actions

#### concurrent programming, 2 forms
##### form 1: communicating sequential processes (CSP)...
- multiple processes running at the same time
- values passed between processes, but variables not shared
- --Go proverb--> `"Don't communicate by sharing memory, share memory by communicating"`
- this is __goroutines & channels__
- major feature of Go. Go aims to bring _simplicity_ to many things, but its approach to concurrent programming is one of its bolder steps at achieving simplicity

##### form 2: shared memory multithreading
- chapter 9

#### goroutines
- a goroutine is a process
- goroutine #1: the main goroutine
    - program starts and calls `main()`
- launch a new goroutine with "go" as in `go myfunction(1, 2, "wazoo")`
    - `myfunction()` starts running, but the line that invoked it _does not wait_ for it to complete
- main goroutine exits ==> all active goroutines abruptly terminated
- lightweight. not threads, though they function the same. do a lot of sleeping and waking
- example: [spinner](https://github.com/adonovan/gopl.io/blob/master/ch8/spinner/main.go)
    - how many goroutines?

#### channels
- channels connect goroutines. short for "communication channel"
- one goroutine puts something into a channel (__SEND__) and another routine takes that something out of the channel (__RECEIVE__)
- a channel is declared to contain a specified TYPE
- zero value: nil
- make an integer channel: `funChannel := make(chan int)`
- SEND x into the channel: `funChannel <- x`
- RECEIVE what's been put into the channel: `foo = <-funChannel`
- CLOSE a channel: `close(funChannel)`
    - check to see if closed `foo, stillOpen = <-funChannel`
    - don't close it twice!
- unbuffered channels...
    - BLOCK ON SEND until item in the channel is received
    - and BLOCK ON RECEIVE until something is loaded into the channel
    - "SYNCHRONOUS" channel... sender cannot continue until something takes my message off the channel... value is received before sender's goroutine reawakes
- pass messages... or signal events
    - value passed (message) may be significant
    - or, adding to the channel may simply signify that it's time to do the next thing, but there is no value being passed, in which case we create a channel of type `struct{}{}` or `bool` or `int` (and ignore the value)
- buffered channels...
    - have a CAPACITY... a number of messages the channel can contain before blocking senders
    - buffered channels do not block on send until the buffer is filled up
    - make a buffered one: `funChannel := make(chan int, 3)`
    - more on using these later...

#### channels as pipelines
[the famous pipelines blog post](https://blog.golang.org/pipelines)

example (pipeline1): generate numbers... square each of these numbers... print each square... uses unbuffered channels
- how many goroutines?

```go
package main

import (
    "fmt"
    "time"
)

const maxNaturals = 5

func main() {
    naturals := make(chan int)
    squares := make(chan int)

    // Number Generator
    go func() {
        for x := 0; x <= maxNaturals; x++ {
            time.Sleep(time.Second)
            naturals <- x
        }
        // we are done producing numbers
        close(naturals)
    }()

    // Squarer
    go func() {
        for x := range naturals { // when naturals closes, this loop exits
            squares <- x * x
        }
        close(squares) // this will shut down our program
    }()

    // Printer (in main goroutine)
    for s := range squares { // when squares closes, this loop exits
        fmt.Println(s)
    }
}

```

#### uni-directional channels
one-way channel, like so: `func foo(mySendChannel chan<- int)` and `func bar(myReceiveChannel <-chan int)`

let's re-do our pipeline example w/uni-directional channels (pipeline-uni)
```go
package main

import (
    "fmt"
    "time"
)

const maxNaturals = 5

func generator(out chan<- int) {
    for x := 0; x <= maxNaturals; x++ {
        time.Sleep(time.Second)
        out <- x
    }
    // we are done producing numbers
    close(out)
}

func squarer(out chan<- int, in <-chan int) {
    for x := range in {
        out <- x * x
    }
    close(out)
}

func printer(in <-chan int) {
    for s := range in {
        fmt.Println(s)
    }
}

func main() {
    naturals := make(chan int)
    squares := make(chan int)

    go generator(naturals)
    go squarer(squares, naturals)
    printer(squares)
}

```

#### pipeline + buffered channels...

we have this...
```go
func main() {
    naturals := make(chan int)
    squares := make(chan int)
    ...
```

but what if our workers work at different rates?

the cake example... BAKER --> ICER --> INSCRIBER
- work at different rates
- very similar to the squarer example, but create channels with **buffers** and have **multiple icers**
- buffered channels so a sender can dump more than one thing into a channel without having to wait... and so a bunch of receivers can grab items off the channel
- example lets you play with channel buffer size and number of workers
- [book code](https://github.com/adonovan/gopl.io/blob/master/ch8/cake/cake.go) (not now)

#### looping in parallel
We want to do the same thing to a bunch of files or whatever. "Concurrency (composition of independently-executing processes) is not parallelism (simultaneous execution of computations, often related)"

This example makes thumbnails for each file received from the channel & returns the number of bytes occupied by the files it creates. It is dummied up to run without external dependencies... just run the main.

It demonstrates the use of a `WaitGroup`

```go
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
```

#### select
- The `select` statement lets us wait on multiple channels (reads &/or writes) and act as soon as one of them is ready. If multiple are ready at the same time, selection is random.

```go
select {
    case <- catChannel:
        etc...
    case myDog := <-dogChannel
        etc...
    case batChannel <- "batman"
        etc...
    default:
        // if nobody is ready when we arrive at this select, do this instead of blocking...
}
```

#### web crawler example

Web Crawler. Terminates when crawling is complete [gist is here](https://gist.github.com/tdarci/89716fbe3947c916b723d3cca6977881):
```go
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
	"sync"
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
	var wg sync.WaitGroup

	// process input
	wg.Add(1) // adding item to crawlResults
	go func() {
		crawlResults <- programArgs() // this gets things rolling...
	}()

	// Notice when we have no more pending operations.
	// At that point, this program is done.
	go func() {
		wg.Wait()
		close(crawlResults) // this will end the program
	}()

	// spawn crawlers (workers that run extract())
	for i := 0; i < concurrentWorkerCount; i++ {
		// each of these makes a crawler to run concurrently. these goroutines die when the program exits
		go func() {
			for link := range unseenLinks { // We will block here until there's something to take off of the channel and keep looping until something closes unseenLinks or main exits
				// We have 2 actions for wg now... +1 for initiating extract() and -1 for taking something off of unseenLinks. So no need to do anything.
				foundLinks := extract(link)
				if foundLinks != nil {
					wg.Add(1) // adding item to crawlResults
					go func() { crawlResults <- foundLinks }()
				}
				wg.Done() // extract() is complete
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
				wg.Add(1) // adding item to unseenLinks
				unseenLinks <- curLink
			}
		}
		wg.Done() // item taken off of crawlResults
	}
}

```

Whew!

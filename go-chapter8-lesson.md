### Chapter 8: Goroutines & Channels

#### concurrent programming, 2 forms
##### form 1: communicating sequential processes (CSP)...
- multiple processes running at the same time-ish
- values passed between processes, but variables not shared --proverb--> "Don't communicate by sharing memory, share memory by communicating"
- this is __goroutines & channels__
- major feature of Go

##### form 2: shared memory multithreading
- chapter 9

#### goroutine mechanics
- a goroutine is a process
- goroutine #1: the main goroutine
    - program starts and calls `main()`
- launch a new goroutine with `go` as in `go myfunction(1, 2, "wazoo")`
    - `myfunction()` starts running, but the line that invoked it does not wait for it to complete
- when main goroutine exits or program is terminated, all active goroutines are abruptly terminated
- these things are lightweight. they are not threads, though they function the same.
- example (spinner): https://github.com/adonovan/gopl.io/blob/master/ch8/spinner/main.go

#### channels
- channels connect goroutines
- one goroutine puts something into a channel (__SEND__) and another routine takes that something out of the channel (__RECEIVE__)
- a channel is declared to contain a specified TYPE
- zero value: nil
- make an integer channel: `funChannel := make(chan int)`
- SEND x through the channel... `funChannel <- x`
- RECEIVE what's been put into the channel... `foo = <-funChannel`
- CLOSE a channel: `close(funChannel)`
    - check to see if closed `foo, stillOpen = <-funChannel`
    - but don't close it twice!
- unbuffered channels...
    - BLOCK ON SEND until item in the channel is received
    - and BLOCK ON RECEIVE until something is loaded into the channel
    - aka: SYNCHRONOUS channel... sender can't continue until something takes my message off the channel... value is received before sender's goroutine reawakes
- channel contents
    - value passed (message) may be significant
    - or, adding to the channel may simply signify that it's time to do the next thing, but there is no value being passed, in which case we create a channel of type `struct{}{}` or `bool` or `int`
- buffered channels...
    - have a CAPACITY... a number of messages the channel can contain before blocking senders
    - buffered channels do not block on send until the buffer is filled up
    - make a buffered one: `funChannel := make(chan int, 3)`
    - more on using these later...

#### channels as pipelines
read this later: https://blog.golang.org/pipelines

example (pipeline1): generate numbers... square each of these numbers... print each square... uses unbuffered channels
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
can declare as uni-directional, like so `func foo(sendChannel chan<- int)` and `func bar(receiveChannel <-chan int)`

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

#### what about buffered channels...

we have this...
```go
func main() {
    naturals := make(chan int)
    squares := make(chan int)
    ...
```

but what if our workers work at different rates?

the cake example... BAKER --> ICER --> INSCRIBER

https://github.com/adonovan/gopl.io/blob/master/ch8/cake/cake.go (not now)

the idea:
- buffered channels so a worker can dump more than one thing into a channel without having to wait
- multiple icers (or bakers or inscribers or whatever)

#### looping in parallel
We want to do the same thing to a bunch of files or whatever. "Concurrency (doing many things at once) is not parallelism (doing the same thing lots of times concurrently)"

example: https://github.com/adonovan/gopl.io/tree/master/ch8/thumbnail (not now)

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

web crawler, parallel processing, but limited to 20 concurrent workers using tokens/semaphores... (skip)
```go
var tokens = make(chan struct{}, 20)

func crawl(url string) []string {
    fmt.Println(url)
    tokens<- struct{}{} // BLOCKS until there is room in the channel's buffer
    list, err := links.Extract(url)
    <-tokens // We're done, so free up a spot in the buffer for a worker
    if err != nil {
        log.Print(err)
    }
    return list
}

func main() {
    worklist := make(chan []string)
    var pendingSendCount int
    pendingSendCount++

    go func(){worklist <- os.Args[1:]}() // process the command line

    seen := make(map[string]bool)
    for ;pendingSendCount > 0; pendingSendCount-- {
        urlList := <-worklist
        for _, curLink := range urlList {
            if !seen[curLink] {
                seen[curLink] = true
                pendingSendCount++
                go func(link string) {
                    worklist<- crawl(link)
                }(curLink)
            }
        }
    }

}
```

- Here's another way to do it... that does not terminate... (skip)
```
package main

import (
	"fmt"
	"log"
	// "os"
)

func Extract(url string) ([]string, error) {
	return nil, nil
}

func programArgs() []string {
	// return os.Args[1:]
	return []string{"foo.bar/wookie"}
}

func crawl(url string) []string {
	fmt.Printf("* Crawling URL: %s\n", url)
	list, err := Extract(url)
	if err != nil {
		log.Printf("Error: %s", err)
	}
	return list
}

func main() {
	worklist := make(chan []string) // channel of lists of urls to process. may have duplicates
	unseenLinks := make(chan string) // channel of de-duped urls


	go func() {worklist <- programArgs()}()

	for i:=0; i < 20; i++ {
		// each of these makes a crawler. 20 run concurrently
		go func() {
			for link := range unseenLinks { // We will block here until there's something to take off of the channel
				// and we keep looping until something closes unseenLinks or main exits
				foundLinks := crawl(link)
				go func() {worklist <- foundLinks}()
			}
		}()
	}

	seen := make(map[string]bool)
	for list := range worklist {
		for _, curLink := range list {
			if !seen[curLink] {
				seen[curLink] = true
				unseenLinks <- curLink
			}
		}
	}
}
```

This one terminates when crawling is complete [gist is here](https://gist.github.com/tdarci/89716fbe3947c916b723d3cca6977881):
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
        pending <- 1
        crawlResults <- programArgs()
    }()

    // Notice when we have no more pending operations.
    // At that point, this program is done.
    go func() {
        // Note that once this function is running it will take items off of pending as soon as they are put there. So no blocking to worry about.
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
            for link := range unseenLinks { // We will block here until there's something to take off of the channel
                // We keep looping until something closes unseenLinks or main exits
                //
                // We have 2 actions for pending channel now... +1 for initiating extract() and -1 for taking something off of unseenLinks. So no need to do anything.
                foundLinks := extract(link)
                if foundLinks != nil {
                    pending <- 1
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
                pending <- 1
                unseenLinks <- curLink
            }
        }
        pending <- -1 // item taken off of crawlResults
    }
}
```

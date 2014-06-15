package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// This file shows the benefits of concurrency and some patterns by attempting
// to download the bodies of several websites and displaying the time it took.
func main() {
	items := []string{
		"http://google.com",
		"http://reddit.com",
		"http://ruby-lang.org",
		"http://golang.org",
	}

	start := time.Now()
	concurrentTimeout(items, 200*time.Millisecond)
	fmt.Printf("Timeout:        \t%.4fs\n\n", time.Since(start).Seconds())
}

type urlLen struct {
	url  string
	size int64
}

func concurrentTimeout(items []string, dt time.Duration) {
	timeout := make(chan bool)
	c := make(chan urlLen)
	defer func() {
		close(timeout)
		close(c)
	}()
	go func() {
		time.Sleep(dt)
		timeout <- true
	}()
	for _, url := range items {
		go countChan(url, c)
	}
	for _ = range items {
		select {
		case r := <-c:
			fmt.Printf("%v\n", r)
		case <-timeout:
			return
		}
	}
}

func countChan(url string, c chan<- urlLen) {
	c <- count(url)
}

func count(url string) urlLen {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	n, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		panic(err)
	}
	return urlLen{url, n}
}

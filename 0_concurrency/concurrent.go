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
	concurrent(items)
	fmt.Printf("Concurrent:     \t%.4fs\n\n", time.Since(start).Seconds())
}

type urlLen struct {
	url  string
	size int64
}

func concurrent(items []string) {
	c := make(chan urlLen)
	defer func() { close(c) }()

	for _, url := range items {
		go countChan(url, c)
	}
	// len(items) == len(reads from channel)
	for _ = range items {
		fmt.Printf("%v\n", <-c)
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

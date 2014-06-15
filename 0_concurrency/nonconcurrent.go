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
	not_concurrent(items)
	fmt.Printf("Not Concurrent: \t%.4fs\n\n", time.Since(start).Seconds())
}

func not_concurrent(items []string) {
	for _, url := range items {
		fmt.Printf("%v %v\n", url, count(url))
	}
}

func count(url string) int64 {
	r, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	n, err := io.Copy(ioutil.Discard, r.Body)
	if err != nil {
		panic(err)
	}
	return n
}

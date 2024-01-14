package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

//type FetchedURL struct {
//	mu        sync.Mutex
//	IsFetched map[string]struct{}
//}

type IsFetched map[string]struct{}

var isFetched = IsFetched{}

//var CheckURL = FetchedURL{
//	IsFetched: map[string]struct{}{},
//}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	fmt.Println("fetching url", url)

	if res, ok := f[url]; ok {
		isFetched[url] = struct{}{}
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func GoTourCrawler(url string, depth int, fetcher Fetcher, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	if depth <= 0 {
		return
	}
	// before fetching check if the url has been already fetched or not.
	// if yes, then dont fetch again
	mu.Lock()
	defer mu.Unlock()
	if _, ok := isFetched[url]; ok {
		fmt.Printf("url: %s has already been fetched\n", url)
		return
	}

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// fmt.Println("fetching for depth", depth-1)
		go GoTourCrawler(u, depth-1, fetcher, wg, mu)
	}

}

func CrawlEx() {
	// main should wait for the urls to be fetched
	var wg sync.WaitGroup
	var mu sync.Mutex

	wg.Add(1)
	GoTourCrawler("https://golang.org/", 4, fetcher, &wg, &mu)
	wg.Wait()
}

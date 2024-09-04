package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
    maxPages           int
}

func (cfg *config) addPageVisit(normalizedURL string) (isFirst bool) {
    cfg.mu.Lock()
    if _, ok := cfg.pages[normalizedURL]; ok {
        cfg.pages[normalizedURL] += 1
        cfg.mu.Unlock()
        return false
    }else {
        cfg.pages[normalizedURL] = 1
        cfg.mu.Unlock()
        return true
    }
}
func (cfg *config) crawlPage(rawCurrentURL string) {
    cfg.concurrencyControl <- struct{}{}
    defer func() {
        <-cfg.concurrencyControl
        cfg.wg.Done()
    }()

    if len(cfg.pages) > cfg.maxPages {
        return
    }

    parsed ,err := url.Parse(rawCurrentURL)
    if err != nil {
        return
    }
    parsedRaw ,err := url.Parse(cfg.baseURL.String())
    if err != nil {
        return
    }

    if parsed.Hostname() != parsedRaw.Hostname() {
        return
    }

    normalizedURL, err := normalizeURL(rawCurrentURL)
    if err != nil {
        return
    }
    
    if cfg.addPageVisit(normalizedURL) == false {
        return
    }

   

    htmlBody, err := getHTML(rawCurrentURL)
    if err != nil {
        fmt.Println(err)
        return
    }
    urls ,err := getURLsFromHTML(htmlBody, cfg.baseURL.String())
    if err != nil {
        return
    }
    
    for _, url := range urls {
        fmt.Printf("found url: %s\nCrawling...", url)
        cfg.wg.Add(1)
        go cfg.crawlPage(url)
    }
}

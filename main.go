package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
)

func getHTML(rawURL string) (string, error) {

    resp, err := http.Get(rawURL)
    if err != nil {
        return "", err
    }
    
    if resp.StatusCode > 400 {
        return "", fmt.Errorf("invalid status code: %d", resp.StatusCode)
    }

    if header := resp.Header.Get("Content-Type");strings.Contains(header, "text/html") == false {
        return "",fmt.Errorf("invalid content type: %s", resp.Header.Get("Content-Type"))
    }

    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    return string(body), nil
}

func main() {
    cmdArgs := os.Args[1:]
    

    if len(cmdArgs) < 3 {
        fmt.Println("not enough arguments")
        os.Exit(1)
    }

    if len(cmdArgs) > 3 {
        fmt.Println("too many arguments")
        os.Exit(1)
    }

    BASE_URL := cmdArgs[0]
    MAX_PAGES,err := strconv.Atoi(cmdArgs[1])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    MAX_CONCURRENCY,err:= strconv.Atoi(cmdArgs[2])
    if err  != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println("starting crawl of: ",BASE_URL)
    parsedURL,err := url.Parse(BASE_URL)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    pages := make(map[string]int) 
    cfg := config{
        pages: pages,
        baseURL: parsedURL,
        mu: &sync.Mutex{},
        concurrencyControl: make(chan struct{},MAX_CONCURRENCY),
        wg: &sync.WaitGroup{},
        maxPages: MAX_PAGES,
    }
    cfg.wg.Add(1)
    go cfg.crawlPage(BASE_URL)
    cfg.wg.Wait()
    printReport(pages, BASE_URL) 
    }


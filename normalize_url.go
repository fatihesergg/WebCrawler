package main

import (
	"net/url"
	"strings"
	"golang.org/x/net/html"
)


func normalizeURL(input string) (string ,error) {
    parsed,err := url.Parse(input)
    if err != nil {
        return "",err
    }
    fullUrl :=  parsed.Hostname() + parsed.Path 
    return fullUrl, nil
}


func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
    urls := []string{}

    doc,err := html.Parse(strings.NewReader(htmlBody))
    if err != nil {
        return []string{},err
    }
    var f func (*html.Node) error 
    f = func(n *html.Node) error{
        if n.Type == html.ElementNode && n.Data == "a" {
            for _, a := range n.Attr {
                if a.Key == "href" {
                    parsed,err := url.Parse(a.Val)
                    if err != nil {
                        return err 
                    }
                    if !parsed.IsAbs() {
                        fullUrl := rawBaseURL + a.Val
                        urls = append(urls, fullUrl)
                    } else {
                        urls = append(urls, a.Val)
                    }
                }
            }
        }
    
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            f(c)
        }
        
        return nil
    }
    
    err =  f(doc)
    if err != nil {
        return []string{},err
    }
    
    return urls,nil
}

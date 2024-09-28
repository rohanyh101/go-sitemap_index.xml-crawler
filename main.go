package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SeoData struct {
	URL             string
	Title           string
	H1              string
	MetaDescription string
	StatusCode      int
}

type DefaultParser struct{}

type Parser interface {
	getSEOData(resp *http.Response) (SeoData, error)
}

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Safari/605.1.15",
	"Mozilla/5.0 (Linux; Android 11; Pixel 5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Mobile Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; WOW64; rv:89.0) Gecko/20100101 Firefox/89.0",
	"Mozilla/5.0 (Linux; Android 10; SM-G950F Build/QP1A.190711.020; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/79.0.3945.136 Mobile Safari/537.36",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:93.0) Gecko/20100101 Firefox/93.0",
	"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
	"Mozilla/5.0 (Linux; U; Android 4.1.1; en-us; Galaxy Nexus Build/JRO03C) AppleWebKit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
	"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	"Mozilla/5.0 (compatible; curl/7.68.0; +https://curl.se/)",
	"curl/7.68.0",
}

func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

func isSitemap(urls []string) ([]string, []string) {
	sitemapFiles := []string{}
	pages := []string{}

	for _, url := range urls {
		if strings.Contains(url, "sitemap") {
			log.Printf("Found sitemap: %s", url)
			sitemapFiles = append(sitemapFiles, url)
		} else {
			pages = append(pages, url)
		}
	}

	return sitemapFiles, pages
}

func extractSiteMapURLs(stratURL string) []string {
	worklist := make(chan []string)
	toCrawl := []string{}
	var n int
	n++

	go func() { worklist <- []string{stratURL} }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, link := range list {
			n++
			go func(link string) {
				response, err := makeRequest(link)
				if err != nil {
					log.Printf("Error making request to %s: %v", link, err)
				}

				urls, err := extractURLs(response)
				if err != nil {
					log.Printf("Error extracting urls from %s: %v", link, err)
				}

				sitemapFiles, pages := isSitemap(urls)
				if sitemapFiles != nil {
					worklist <- sitemapFiles
				}

				for _, page := range pages {
					toCrawl = append(toCrawl, page)
				}
			}(link)
		}
	}

	return toCrawl
}

func makeRequest(url string) (*http.Response, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", randomUserAgent())
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func scrapeURLs(urls []string, parser Parser, concurrency int) []SeoData {
	tokens := make(chan struct{}, concurrency)
	worklist := make(chan []string)
	results := []SeoData{}
	var n int
	n++

	go func() { worklist <- urls }()
	for ; n > 0; n-- {
		list := <-worklist
		for _, url := range list {
			if url != "" {
				n++
				go func(url string, tokens chan struct{}) {
					log.Printf("Scraping URL: %s", url)
					res, err := scrapePage(url, tokens, parser)
					if err != nil {
						log.Printf("Error scraping URL: %s, %v", url, err)
						worklist <- []string{}
					}

					results = append(results, res)
					worklist <- []string{}
				}(url, tokens)
			}
		}
	}

	return results
}

func extractURLs(response *http.Response) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	urls := []string{}
	// doc.Find("a").Each(func(i int, s *goquery.Selection) {
	// 	url, _ := s.Attr("href")
	// 	urls = append(urls, url)
	// })

	s := doc.Find("loc")
	for i := range s.Nodes {
		url := s.Eq(i).Text()
		urls = append(urls, url)
	}

	return urls, nil
}

func scrapePage(url string, token chan struct{}, parser Parser) (SeoData, error) {
	res, err := crawlPage(url, token)
	if err != nil {
		return SeoData{}, err
	}

	data, err := parser.getSEOData(res)
	if err != nil {
		return SeoData{}, err
	}

	return data, nil
}

func crawlPage(url string, tokens chan struct{}) (*http.Response, error) {
	tokens <- struct{}{}
	defer func() { <-tokens }()

	res, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (d DefaultParser) getSEOData(resp *http.Response) (SeoData, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return SeoData{}, err
	}

	result := SeoData{}
	result.URL = resp.Request.URL.String()
	result.StatusCode = resp.StatusCode
	result.Title = doc.Find("title").Text()
	result.H1 = doc.Find("h1").Text()
	result.MetaDescription, _ = doc.Find("meta[name=description]").Attr("content")

	return result, nil
}

func ScrapeSiteMap(url string, parser Parser, concurrency int) []SeoData {
	results := extractSiteMapURLs(url)
	res := scrapeURLs(results, parser, concurrency)
	return res
}

func main() {
	p := DefaultParser{}
	results := ScrapeSiteMap("https://yoast.com/sitemap_index.xml", p, 10)
	for _, res := range results {
		fmt.Printf("%+v\n", res)
	}
}

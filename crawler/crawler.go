package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"simple-go-crawler/common"
	"simple-go-crawler/constant"
	"simple-go-crawler/scraper"
	"strings"
	"sync"
)

type Crawler struct {
	baseURL        *url.URL
	productCounter *scraper.ProductCounter
	visitTracker   *VisitTracker
	urlQueue       *URLQueue
	wg             sync.WaitGroup
}

func NewCrawler(baseURL string) (*Crawler, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	return &Crawler{
		baseURL:        parsedURL,
		productCounter: scraper.NewProductCounter(),
		visitTracker:   NewVisitedURLTracker(),
		urlQueue:       NewURLQueue(),
	}, nil
}

func (c *Crawler) Start() {
	c.urlQueue.Add(c.baseURL.String())
	//c.Crawl()
}

//func (c *Crawler) Crawl() {
//	for {
//		if c.productCounter.GetCount() == 300 {
//			break
//		}
//		if c.urlQueue.IsEmpty() {
//			fmt.Println("Queue is empty, waiting for new URL...")
//			continue
//		}
//
//		urlsToProcess := c.urlQueue.GetQueue()
//		//c.urlQueue.Reset()
//		for k := range urlsToProcess {
//			c.urlQueue.Remove(k)
//			if c.productCounter.GetCount() == 300 {
//				break
//			}
//			//c.wg.Add(1) // Increment the WaitGroup counter
//			//go func(currentURL string) {
//			//	defer c.wg.Done() // Decrement the WaitGroup counter when this goroutine finishes
//			//	fmt.Printf("Processing URL: %s\n", currentURL)
//			//	htmlBody, err := common.GetWebPage(currentURL)
//			//	if err != nil {
//			//		fmt.Printf("Error getting web page %s: %v\n", currentURL, err)
//			//		return
//			//	}
//			//	// extracts the links from the html body and possibly adds new URLs to c.urlQueue
//			//	c.extractData(htmlBody)
//			//	fmt.Printf("Finished processing %s. Found %d links on queue.\n", currentURL, c.urlQueue.GetLength())
//			//}(k)
//			fmt.Printf("Processing URL: %s\n", k)
//			htmlBody, err := common.GetWebPage(k)
//			if err != nil {
//				fmt.Printf("Error getting web page %s: %v\n", k, err)
//				return
//			}
//			// extracts the links from the html body and possibly adds new URLs to c.urlQueue
//			c.extractData(htmlBody)
//			fmt.Printf("Finished processing %s. Found %d links on queue.\n", k, c.urlQueue.GetLength())
//		}
//
//		// Wait for all goroutines launched in this iteration to finish
//		// This is crucial to prevent premature exit and ensure all data is extracted
//		//c.wg.Wait()
//		fmt.Println("All goroutines for this batch completed.")
//		fmt.Println("Total Products: ", c.productCounter.GetCount())
//		// Small sleep to avoid busy-looping if there are always URLs
//		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
//	}
//}

func (c *Crawler) Crawl(urlStr string) error {
	htmlBody, err := common.GetWebPage(urlStr)
	if err != nil {
		return err
	}

	// extracts the links from the html body
	c.extractData(htmlBody)

	fmt.Printf("Found %d links on %s worth visiting \n", c.urlQueue.GetLength(), urlStr)
	//time.Sleep(1 * time.Second)
	return nil
}

// Extracts links from an HTML node recursively
func (c *Crawler) extractData(n *html.Node) {
	if n == nil {
		return
	}

	// extracts All URLs
	if n.Type == html.ElementNode && n.Data == "a" {
		isProductCard := false
		for _, attr := range n.Attr {
			// <a data-testid="product-card-image-link" href="/products/IQ1401">
			if attr.Key == "data-testid" && attr.Val == constant.PRODUCT_URL_IDENTIFIER || attr.Key == "class" && attr.Val == "_product-card__link_o6rgp_73" {
				isProductCard = true
			}
			if attr.Key == "href" {
				resolvedURL := c.resolveURL(attr.Val)
				if resolvedURL != "" {
					c.urlQueue.Add(resolvedURL)
					fmt.Println("NEW URL: " + resolvedURL)
					if isProductCard {
						scraper.ExtractProductData(resolvedURL, c.productCounter)
					}
				}
			}
		}
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		c.extractData(child)
	}
}

// converts relative URLs to absolute URLs, returns the original URL if it's already absolute
func (c *Crawler) resolveURL(href string) string {
	// Check if it's alreday an absolute URL (starts with http:// or https://)
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	relURL, err := url.Parse(href)
	if err != nil {
		return ""
	}
	absURL := c.baseURL.ResolveReference(relURL)
	return absURL.String()
}

// checks if the URL should be crawled (same domain)
func (c *Crawler) shouldCrawl(urlStr string) bool {
	return strings.Contains(urlStr, "https://www.adidas.jp/")
}

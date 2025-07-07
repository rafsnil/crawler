package crawler

import (
	"fmt"
	"golang.org/x/net/html"
	"log"
	"math/rand"
	"net/url"
	"simple-go-crawler/common"
	"simple-go-crawler/constant"
	"simple-go-crawler/scraper"
	"strings"
	"sync"
	"time"
)

type Crawler struct {
	baseURL        *url.URL
	ProductCounter *scraper.ProductCounter
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
		ProductCounter: scraper.NewProductCounter(),
		visitTracker:   NewVisitedURLTracker(),
		urlQueue:       NewURLQueue(),
	}, nil
}

var LinksLoaded bool = false

func (c *Crawler) Start() {
	c.urlQueue.Add(c.baseURL.String())
	c.Crawl()
}

func (c *Crawler) Crawl() {
	const emptyQueueTimeout = 30 * time.Second // 30secs
	lastNonEmptyTime := time.Now()
	retry := 1
	for {
		if c.ProductCounter.GetCount() >= constant.DATA_LIMIT {
			break
		}

		//if c.ProductCounter.GetCount() == 40 {
		//	time.Sleep(100 * time.Second)
		//	fmt.Println("Sleeping for 100 seconds to avoid getting blocked")
		//}

		if c.urlQueue.IsEmpty() {
			if time.Since(lastNonEmptyTime) > emptyQueueTimeout {
				fmt.Println("Queue has been empty for too long. Exiting crawl.")
				break
			}
			time.Sleep(2 * time.Second) // avoid tight spinning
			continue
		} else {
			lastNonEmptyTime = time.Now()
		}

		urlsToProcess := c.urlQueue.GetQueue()

		for k := range urlsToProcess {
			c.urlQueue.Remove(k)

			if c.ProductCounter.GetCount() >= constant.DATA_LIMIT {
				break
			}

			fmt.Printf("Processing URL: %s\n", k)
			stringBody, htmlBody, err := common.GetWebPage(k)
			if err != nil {
				fmt.Printf("Error getting web page %s: %v\n", k, err)
				continue
			}

			allGood := checkIfHtmlResponseIsOkay(stringBody)
			if !allGood {
				if retry == 10 {
					break
				}
				c.urlQueue.Add(k)
				//fmt.Println(stringBody)
				fmt.Println("Got blocked by the website")
				fmt.Println("Sleeping for 10 seconds before retrying...")
				time.Sleep(10 * time.Second)
				retry++
				continue
			}

			c.extractData(htmlBody)

			fmt.Printf("Finished processing %s. Found %d links on queue.\n", k, c.urlQueue.GetLength())
			fmt.Println("Total Products: ", c.ProductCounter.GetCount())
			d := time.Duration(rand.Intn(5)+3) * time.Second
			// save visited url to file
			//err = common.SaveURLToFile(k, constant.VISITED)
			//if err != nil {
			//	fmt.Printf("Error saving URL %s: %v\n", k, err)
			//}
			fmt.Println("sleeping for", d, "before next request to load web page")
			time.Sleep(d)
		}

	}
}

func checkIfHtmlResponseIsOkay(body string) bool {
	return strings.Contains(body, constant.SCRIPT_DATA_IDENTIFIER)
}

// Extracts links from an HTML node recursively
func (c *Crawler) extractData(n *html.Node) {
	scriptData, err := common.GetWebPageScriptData(n)
	if err != nil {
		fmt.Printf("Error extracting script data: %v\n", err)
		return
	}
	if scriptData.Props.PageProps.PageType == constant.PAGE_TYPE_PRODUCT_LISTING {
		for _, product := range scriptData.Props.PageProps.Products {
			fmt.Printf("Product found: %s\n", product.Id)
			if c.ProductCounter.GetCount() >= constant.DATA_LIMIT {
				fmt.Println("Product limit reached, stopping extraction.")
				return
			}
			scraper.ExtractProductData(product.Id, c.ProductCounter)
		}
	} else if scriptData.Props.PageProps.PageType == constant.PAGE_TYPE_LANDING {
		for _, layout := range scriptData.Props.PageProps.Layouts {
			for _, content := range layout.Contents {
				for _, cta := range content.CTAs {
					if cta.RelativeUrl != "" {
						resolvedURL := c.resolveURL(cta.RelativeUrl)
						if resolvedURL != "" && c.shouldCrawl(resolvedURL) {
							c.urlQueue.Add(resolvedURL)
							//_ = common.SaveURLToFile(resolvedURL, constant.QUEUE)
						}
					}
				}
			}
		}
	}
	if len(scriptData.Props.PageProps.Products) > 0 {
		for _, product := range scriptData.Props.PageProps.Products {
			fmt.Printf("Product found: %s\n", product.Id)
			if c.ProductCounter.GetCount() >= constant.DATA_LIMIT {
				fmt.Println("Product limit reached, stopping extraction.")
				return
			}
			scraper.ExtractProductData(product.Id, c.ProductCounter)
		}
	}

	if scriptData.Props.PageProps.PageSeo != nil {
		if len(scriptData.Props.PageProps.PageSeo.PlpLinks) > 0 {
			for _, relUrl := range scriptData.Props.PageProps.PageSeo.PlpLinks {
				if relUrl.Url != "" {
					resolvedURL := c.resolveURL(relUrl.Url)
					if resolvedURL != "" && c.shouldCrawl(resolvedURL) {
						c.urlQueue.Add(resolvedURL)
					}
				}
			}
		}
	}

	if scriptData.Props.PageProps.NavigationData != nil {
		if len(scriptData.Props.PageProps.NavigationData.NavItems) > 0 {
			for _, navItem := range scriptData.Props.PageProps.NavigationData.NavItems {
				if navItem.Href != "" {
					resolvedURL := c.resolveURL(navItem.Href)
					if resolvedURL != "" && c.shouldCrawl(resolvedURL) {
						c.urlQueue.Add(resolvedURL)
					}
				}
			}
		}
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
	rawURL := "https://www.adidas.jp/"
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}

	absURL := parsedURL.ResolveReference(relURL)
	return absURL.String()
}

// checks if the URL should be crawled (same domain)
func (c *Crawler) shouldCrawl(urlStr string) bool {
	return strings.Contains(urlStr, "https://www.adidas.jp/") || strings.Contains(urlStr, "https://shop.adidas.jp/")
}

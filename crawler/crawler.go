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
	c.Crawl()
}

func (c *Crawler) Crawl() {
	const emptyQueueTimeout = 30 * time.Second // 30secs
	lastNonEmptyTime := time.Now()
	retry := 1
	for {
		if c.productCounter.GetCount() >= 300 {
			break
		}

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

			if c.productCounter.GetCount() >= 200 {
				break
			}

			fmt.Printf("Processing URL: %s\n", k)
			stringBody, htmlBody, err := common.GetPageHTMLByPassingBot(k)
			if err != nil {
				fmt.Printf("Error getting web page %s: %v\n", k, err)
				continue
			}

			allGood := checkIfHtmlResponseIsOkay(stringBody)
			if !allGood {
				c.urlQueue.Add(k)
				if retry == 10 {
					break
				}
				fmt.Println("Got blocked by the website")
				fmt.Println("Trying with bot bypassing method...")
				htmlStr, htmlContent, err := common.GetPageHTMLByPassingBot(k)
				if err != nil {
					fmt.Printf("Error getting web page %s: %v\n", k, err)
					break
				}
				alright := checkIfHtmlResponseIsOkay(htmlStr)
				if alright {
					c.urlQueue.Remove(k)
					c.extractData(htmlContent)
				}

				fmt.Println("Sleeping for 30 seconds before retrying...")
				time.Sleep(30 * time.Second)
				retry++
				continue
			}

			c.extractData(htmlBody)

			fmt.Printf("Finished processing %s. Found %d links on queue.\n", k, c.urlQueue.GetLength())
			fmt.Println("Total Products: ", c.productCounter.GetCount())
		}
		time.Sleep(time.Duration(rand.Intn(10)+5) * time.Second)
	}
}

func checkIfHtmlResponseIsOkay(body string) bool {
	return strings.Contains(body, constant.SCRIPT_DATA_IDENTIFIER)
}

//func (c *Crawler) Crawl(urlStr string) error {
//	htmlBody, err := common.GetWebPage(urlStr)
//	if err != nil {
//		return err
//	}
//
//	// extracts the links from the html body
//	c.extractData(htmlBody)
//
//	//fmt.Printf("Found %d links on %s worth visiting \n", c.urlQueue.GetLength(), urlStr)
//	//time.Sleep(1 * time.Second)
//	return nil
//}

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
			scraper.ExtractProductData(product.Id, c.productCounter)
		}
	} else if scriptData.Props.PageProps.PageType == constant.PAGE_TYPE_LANDING {
		for _, layout := range scriptData.Props.PageProps.Layouts {
			for _, content := range layout.Contents {
				for _, cta := range content.CTAs {
					if cta.RelativeUrl != "" {
						resolvedURL := c.resolveURL(cta.RelativeUrl)
						if resolvedURL != "" && c.shouldCrawl(resolvedURL) {
							c.urlQueue.Add(resolvedURL)
						}
					}
				}
			}
		}
	}
	if len(scriptData.Props.PageProps.Products) > 0 {
		for _, product := range scriptData.Props.PageProps.Products {
			fmt.Printf("Product found: %s\n", product.Id)
			scraper.ExtractProductData(product.Id, c.productCounter)
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

//func (c *Crawler) extractUrl(n *html.Node) {
//	if n == nil {
//		return
//	}
//
//	// extracts All URLs
//	if n.Type == html.ElementNode && n.Data == "a" {
//		for _, attr := range n.Attr {
//			// <a data-testid="product-card-image-link" href="/products/IQ1401">
//			if attr.Key == "data-testid" && attr.Val == constant.PRODUCT_URL_IDENTIFIER || attr.Key == "class" && attr.Val == "_product-card__link_o6rgp_73" {
//				isProductCard = true
//			}
//			if attr.Key == "href" {
//				resolvedURL := c.resolveURL(attr.Val)
//				if resolvedURL != "" {
//					c.urlQueue.Add(resolvedURL)
//				}
//			}
//		}
//	}
//	for child := n.FirstChild; child != nil; child = child.NextSibling {
//		c.extractData(child)
//	}
//}

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

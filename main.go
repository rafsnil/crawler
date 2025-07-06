package main

import (
	"fmt"
	"log"
	"simple-go-crawler/common"

	"github.com/gocolly/colly"
)

func main() {
	// Create a new collector
	c := colly.NewCollector(
		colly.AllowedDomains("shop.adidas.jp", "www.adidas.jp"),
	)

	// Set a custom User-Agent to avoid 403 errors
	// Set headers to mimic a real browser
	c.OnRequest(func(r *colly.Request) {
		common.SetCommonHeaders(r)
	})
	// On every <a> tag with href, print the link
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link != "" {
			fmt.Println(link)
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Error: %v\nStatus Code: %d\n", err, r.StatusCode)
	})

	// Start scraping
	err := c.Visit("https://shop.adidas.jp/men/")
	if err != nil {
		log.Fatal(err)
	}
}

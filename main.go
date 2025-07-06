package main

import (
	"fmt"
	"simple-go-crawler/crawler"
	"simple-go-crawler/scraper"
	"time"
)

func main() {
	startURL := "https://shop.adidas.jp/men"

	crawlr, err := crawler.NewCrawler(startURL)
	if err != nil {
		fmt.Printf("Error creating crawler: %v\n", err)
		return
	}

	// Start time
	start := time.Now()

	// Run the crawler
	crawlr.Start()

	// End time
	elapsed := time.Since(start)
	fmt.Printf("⏱️ Crawling completed in %s\n", elapsed)
	scraper.GlobalOutput.PrintToExcelSheet()
}

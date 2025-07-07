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
	done := make(chan struct{})
	go func() {
		for {
			fmt.Printf("\n\n\n📦 Product count so far: %d\n\n\n", len(scraper.GlobalOutput.Products))
			time.Sleep(10 * time.Second)
		}
	}()

	// Start time
	start := time.Now()

	// Run the crawler
	crawlr.Start()
	close(done)
	// End time
	elapsed := time.Since(start)

	fmt.Printf("⏱️ Crawling completed in %s\n", elapsed)

	scraper.GlobalOutput.PrintToExcelSheet()
}

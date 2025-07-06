package main

//
//import (
//	"fmt"
//	"net/url"
//	"simple-go-crawler/common"
//	"sync"
//	"time"
//)
//
//// Dummy implementations for demonstration purposes
//type ProductCounter struct {
//	count int
//	mu    sync.Mutex
//}
//
//func (pc *ProductCounter) GetCount() int {
//	pc.mu.Lock()
//	defer pc.mu.Unlock()
//	return pc.count
//}
//
//func (pc *ProductCounter) Increment() {
//	pc.mu.Lock()
//	defer pc.mu.Unlock()
//	pc.count++
//}
//
//type VisitTracker struct{} // Placeholder
//
//type URLQueue struct {
//	queue map[string]bool
//	mu    sync.Mutex
//}
//
//func NewURLQueue() *URLQueue {
//	return &URLQueue{
//		queue: make(map[string]bool),
//	}
//}
//
//func (uq *URLQueue) Add(url string) {
//	uq.mu.Lock()
//	defer uq.mu.Unlock()
//	uq.queue[url] = true
//}
//
//func (uq *URLQueue) IsEmpty() bool {
//	uq.mu.Lock()
//	defer uq.mu.Unlock()
//	return len(uq.queue) == 0
//}
//
//func (uq *URLQueue) GetQueue() map[string]bool {
//	uq.mu.Lock()
//	defer uq.mu.Unlock()
//	// Return a copy to prevent external modification
//	copiedQueue := make(map[string]bool)
//	for k, v := range uq.queue {
//		copiedQueue[k] = v
//	}
//	return copiedQueue
//}
//
//func (uq *URLQueue) GetLength() int {
//	uq.mu.Lock()
//	defer uq.mu.Unlock()
//	return len(uq.queue)
//}
//
//func (uq *URLQueue) Remove(url string) {
//	uq.mu.Lock()
//	defer uq.mu.Unlock()
//	delete(uq.queue, url)
//}
//
//// Crawler struct as provided
//type Crawler struct {
//	baseURL        *url.URL
//	productCounter *ProductCounter // Changed to ProductCounter for simplicity
//	visitTracker   *VisitTracker
//	urlQueue       *URLQueue
//	wg             sync.WaitGroup
//}
//
//// NewCrawler creates a new Crawler instance
//func NewCrawler() *Crawler {
//	baseURL, _ := url.Parse("http://example.com") // Dummy base URL
//	return &Crawler{
//		baseURL:        baseURL,
//		productCounter: &ProductCounter{},
//		visitTracker:   &VisitTracker{},
//		urlQueue:       NewURLQueue(),
//	}
//}
//
//func (c *Crawler) Crawl() {
//	// Initial URLs for demonstration
//	c.urlQueue.Add("http://example.com/page1")
//	c.urlQueue.Add("http://example.com/page2")
//	c.urlQueue.Add("http://example.com/page3")
//
//	for {
//		// Condition to stop crawling
//		if c.productCounter.GetCount() >= 300 {
//			fmt.Println("Product count reached 300. Stopping crawl.")
//			break
//		}
//
//		if c.urlQueue.IsEmpty() {
//			fmt.Println("Queue is empty, waiting for new URL...")
//			// You might want to add a mechanism to truly stop if queue is empty
//			// and no new URLs are expected, or if a certain timeout occurs.
//			// For now, it will keep looping and checking.
//			time.Sleep(1 * time.Second) // Prevent busy-waiting
//			continue
//		}
//
//		// Get a copy of the current queue to iterate over, then clear the original
//		// This prevents processing the same URLs multiple times in a single iteration
//		// and allows new URLs to be added to the queue while processing old ones.
//		urlsToProcess := c.urlQueue.GetQueue()
//		for k := range urlsToProcess {
//			c.urlQueue.Remove(k) // Remove URL from the queue once it's picked for processing
//
//			c.wg.Add(1) // Increment the WaitGroup counter
//			go func(currentURL string) {
//				defer c.wg.Done() // Decrement the WaitGroup counter when this goroutine finishes
//
//				fmt.Printf("Processing URL: %s\n", currentURL)
//
//				htmlBody, err := common.GetWebPage(currentURL) // Assuming common.GetWebPage exists
//				if err != nil {
//					fmt.Printf("Error getting web page %s: %v\n", currentURL, err)
//					return // Important: return from this goroutine on error
//				}
//				// extracts the links from the html body and possibly adds new URLs to c.urlQueue
//				c.extractData(htmlBody) // This method needs to handle adding new URLs to c.urlQueue
//				fmt.Printf("Finished processing %s. Found %d links on queue.\n", currentURL, c.urlQueue.GetLength())
//			}(k) // Pass k as an argument to the goroutine to avoid closure issues
//		}
//
//		// Wait for all goroutines launched in this iteration to finish
//		// This is crucial to prevent premature exit and ensure all data is extracted
//		c.wg.Wait()
//		fmt.Println("All goroutines for this batch completed.")
//
//		// Small sleep to avoid busy-looping if there are always URLs
//		time.Sleep(100 * time.Millisecond)
//	}
//	fmt.Println("Crawler finished.")
//}
//
//// extractData is a placeholder for your actual data extraction logic.
//// It should parse the htmlBody, extract relevant information, and potentially
//// add new URLs to c.urlQueue.
//func (c *Crawler) extractData(htmlBody string) {
//	// Simulate some work
//	time.Sleep(50 * time.Millisecond)
//
//	// Simulate finding new URLs and adding them to the queue
//	// In a real scenario, you would parse htmlBody for actual links
//	if c.productCounter.GetCount() < 300 { // Only add if we still need products
//		c.productCounter.Increment()
//		fmt.Printf("Product count: %d\n", c.productCounter.GetCount())
//		// Add some dummy URLs for next iteration
//		if c.productCounter.GetCount()%5 == 0 { // Add new URLs every 5 products
//			c.urlQueue.Add(fmt.Sprintf("http://example.com/new_page_%d", c.productCounter.GetCount()))
//		}
//	}
//}
//
//func main() {
//	crawler := NewCrawler()
//	crawler.Crawl()
//}

package common

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net/http"
	"simple-go-crawler/constant"
	"simple-go-crawler/dto"
	"time"
)

func GetWebPageScriptData(body *html.Node) (*dto.WebPageScriptData, error) {
	var result dto.WebPageScriptData
	var findScript func(*html.Node) *html.Node
	findScript = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode && n.Data == "script" {
			for _, attr := range n.Attr {
				if attr.Key == "id" && attr.Val == constant.SCRIPT_DATA_IDENTIFIER {
					return n
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if res := findScript(c); res != nil {
				return res
			}
		}
		return nil
	}

	scriptNode := findScript(body)
	if scriptNode == nil {
		return nil, errors.New("script tag with data-mf-id not found")
	}

	if scriptNode.FirstChild == nil {
		return nil, errors.New("script tag found but no content")
	}

	jsonData := scriptNode.FirstChild.Data
	//fmt.Printf("%s", jsonData)

	decoder := json.NewDecoder(bytes.NewBufferString(jsonData))
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetWebPage(url string) (string, *html.Node, error) {
	// Create browser options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("incognito", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var htmlContent string
	tasks := chromedp.Tasks{
		network.Enable(),
		chromedp.Navigate(url),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(5 * time.Second), // Wait for challenges
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	}

	if err := chromedp.Run(ctx, tasks); err != nil {
		return "", nil, fmt.Errorf("chromedp run failed: %v", err)
	}

	_, node, err := GetParsedHTML(bytes.NewReader([]byte(htmlContent)))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return htmlContent, node, nil
}

//func GetWebPage(url string) (string, *html.Node, error) {
//	req, err := http.NewRequest(http.MethodGet, url, nil)
//	if err != nil {
//		return "", nil, err
//	}
//
//	SetCommonHeaders(req, false)
//
//	jar, _ := cookiejar.New(nil)
//	client := &http.Client{
//		Jar:     jar,
//		Timeout: 10 * time.Second,
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		return "", nil, fmt.Errorf("request failed: %v", err)
//	}
//	defer resp.Body.Close()
//
//	if resp.StatusCode != http.StatusOK {
//		return "", nil, fmt.Errorf("received status code %d for url: %s", resp.StatusCode, url)
//	}
//
//	// Read the response body
//	bodyBytes, err := io.ReadAll(resp.Body)
//	if err != nil {
//		return "", nil, fmt.Errorf("failed to read response body: %v", err)
//	}
//
//	// Print response body as string
//	//stringBody := string(bodyBytes)
//	//fmt.Printf("Response body: %s\n", string(bodyBytes))
//	return GetParsedHTML(bytes.NewReader(bodyBytes))
//}

// MakeAPICall makes a GET request to the specified URL and returns the response body as a string.
func MakeAPICall(url string) ([]byte, error) {
	time.Sleep(time.Duration(rand.Intn(3)+3) * time.Second)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	SetCommonHeaders(req, true)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d for url: %s", resp.StatusCode, url)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// check if gzip
	if bytes.HasPrefix(body, []byte{0x1f, 0x8b}) {
		body, err = convertGzip(body)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress content: %v", err)
		}
	}
	return body, nil
}

var referers = []string{
	"https://www.adidas.jp/",
	"https://www.google.com/",
	"https://www.adidas.jp/men",
}
var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
}

func SetCommonHeaders(req *http.Request, apiCall bool) {
	rand.Seed(time.Now().UnixNano())
	ref := referers[rand.Intn(len(referers))]
	ua := userAgents[rand.Intn(len(userAgents))]
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Host", "www.adidas.jp")

	if apiCall {
		req.Header.Set("Referer", "https://www.adidas.jp/")
	} else {
		req.Header.Set("Referer", ref)
	}
}

//var (
//	userAgents = []string{
//		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
//		//"Mozilla/5.0 (Macintosh; Intel Mac OS X 13_4_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.5735.198 Safari/537.36",
//		//"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.5993.88 Safari/537.36",
//		//"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:118.0) Gecko/20100101 Firefox/118.0",
//	}
//
//	acceptLanguages = []string{
//		"en-US,en;q=0.9,ja;q=0.9",
//	}
//
//	acceptHeaders = []string{
//		"text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8",
//	}
//
//	referers = []string{
//		"https://www.adidas.jp/",
//		"https://www.google.com/",
//		"https://www.adidas.jp/men",
//		"https://www.adidas.jp/women",
//		"https://www.adidas.jp/sports",
//	}
//
//	secChUAs = []string{
//		"\"Not.A/Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"",
//	}
//
//	secChUaPlatforms = []string{
//		"\"Windows\"",
//		//"\"macOS\"",
//		//"\"Linux\"",
//	}
//)
//
//func SetCommonHeaders(req *http.Request, apiCall bool) {
//	rand.Seed(time.Now().UnixNano())
//
//	req.Header.Set("User-Agent", userAgents[rand.Intn(len(userAgents))])
//	req.Header.Set("Accept-Language", acceptLanguages[rand.Intn(len(acceptLanguages))])
//	req.Header.Set("Accept", acceptHeaders[rand.Intn(len(acceptHeaders))])
//	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
//	req.Header.Set("Sec-Ch-Ua", secChUAs[rand.Intn(len(secChUAs))])
//	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
//	req.Header.Set("Sec-Ch-Ua-Platform", secChUaPlatforms[rand.Intn(len(secChUaPlatforms))])
//	req.Header.Set("Upgrade-Insecure-Requests", "1")
//	req.Header.Set("Cache-Control", "no-cache")
//	req.Header.Set("Pragma", "no-cache")
//	req.Header.Set("Connection", "keep-alive")
//
//	// Set Referer
//	if apiCall {
//		req.Header.Set("Referer", "https://www.adidas.jp/")
//	} else {
//		req.Header.Set("Referer", referers[rand.Intn(len(referers))])
//	}
//}

func GetParsedHTML(body io.Reader) (string, *html.Node, error) {
	// Read the entire body
	content, err := io.ReadAll(body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read body: %v", err)
	}

	// Check if content is gzipped
	// The gzip format always begins with the two bytes 0x1f and 0x8b
	if bytes.HasPrefix(content, []byte{0x1f, 0x8b}) {
		content, err = convertGzip(content)
		if err != nil {
			return "", nil, fmt.Errorf("failed to decompress content: %v", err)
		}
	}

	// Print preview of decompressed content
	//fmt.Printf("Response body preview: %s\n", string(content[:min(len(content), 500)]))

	/*
		| Field         | Description                                                      |
		| ------------- | ---------------------------------------------------------------- |
		| `Type`        | Tells you what kind of node it is (e.g., Element, Text, Comment) |
		| `Data`        | Tag name like `"div"`, `"a"`, `"ul"`, etc.                       |
		| `Attr`        | List of attributes (like `href`, `class`)                        |
		| `FirstChild`  | First nested tag inside it                                       |
		| `NextSibling` | The next tag at the same level                                   |
	*/
	// Parse the HTML content
	//fmt.Printf(string(content))
	doc, err := html.Parse(bytes.NewReader(content))
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse HTML: %v", err)
	}
	return string(content), doc, nil
}

func convertGzip(content []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

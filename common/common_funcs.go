package common

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net/http"
	"time"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Safari/605.1.15",
	//"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 17_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.0 Mobile/15E148 Safari/604.1",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:120.0) Gecko/20100101 Firefox/120.0",
}

func GetWebPage(url string) (*html.Node, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	//SetCommonHeaders(req)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received status code %d for url: %s", resp.StatusCode, url)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Print response body as string
	fmt.Printf("Response body: %s\n", string(bodyBytes))

	return GetParsedHTML(bytes.NewReader(bodyBytes))
}

// MakeAPICall makes a GET request to the specified URL and returns the response body as a string.
func MakeAPICall(url string) ([]byte, error) {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	//req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	//req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.9") // can't comment this out, as this causes timeout issues
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	//req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	//req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	////req.Header.Set("Upgrade-Insecure-Requests", "1")
	//req.Header.Set("Referer", "https://www.adidas.jp/")
	//SetCommonHeaders(req)

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

//func SetCommonHeaders(req *http.Request) {
//	//NOTE FOR SELF: The first attempt to 'GET' the base url caused a timeout,
//	// then added some common headers to mimic a browser request, however, that too generated a 403
//	// after some trial and error, I found that the following headers are enough to get a 200 OK response and the html content
//	// Enhanced headers to better mimic a modern browser
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
//	//req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
//	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
//	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.9") // can't comment this out, as this causes timeout issues
//	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
//	//req.Header.Set("Cache-Control", "max-age=0") // does not cause any issue if commented out
//	/***
//	These headers are [Client Hints](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers#client_hints) used by browsers to provide additional information to the server about the client:
//
//	- `Sec-Ch-Ua`: Identifies the browser brand and version (e.g., Chromium, Google Chrome).
//	- `Sec-Ch-Ua-Mobile`: Indicates if the client is on a mobile device (`?0` means not mobile).
//	- `Sec-Ch-Ua-Platform`: Specifies the platform or operating system (e.g., "Windows").
//
//	They help servers deliver optimized content or features based on the client’s environment.
//	/
//	*/
//	req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
//	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
//	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
//	req.Header.Set("Referer", "https://www.google.com/") // Needed to make API call
//
//}

func SetCommonHeaders(req *colly.Request) {
	rand.Seed(time.Now().UnixNano())
	_ = userAgents[rand.Intn(len(userAgents))]

	req.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Headers.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.9")
	req.Headers.Set("Accept-Encoding", "gzip, deflate, br")
	req.Headers.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Headers.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Headers.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	req.Headers.Set("Referer", "https://www.google.com/")
}

func GetParsedHTML(body io.Reader) (*html.Node, error) {
	// Read the entire body
	content, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}

	// Check if content is gzipped
	// The gzip format always begins with the two bytes 0x1f and 0x8b
	if bytes.HasPrefix(content, []byte{0x1f, 0x8b}) {
		content, err = convertGzip(content)
		if err != nil {
			return nil, fmt.Errorf("failed to decompress content: %v", err)
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
	return html.Parse(bytes.NewReader(content))
}

func convertGzip(content []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

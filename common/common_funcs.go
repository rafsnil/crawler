package common

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/stealth"
	"golang.org/x/net/html"
	"io"
	"math/rand"
	"net/http"
	"simple-go-crawler/constant"
	"simple-go-crawler/dto"
	"strings"
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
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", nil, err
	}

	SetCommonHeaders(req, false)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("received status code %d for url: %s", resp.StatusCode, url)
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Print response body as string
	//stringBody := string(bodyBytes)
	//fmt.Printf("Response body: %s\n", string(bodyBytes))
	output, err := GetParsedHTML(bytes.NewReader(bodyBytes))
	if err != nil {
		return "", nil, fmt.Errorf("failed to read response body: %v", err)
	}
	return string(bodyBytes), output, nil
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

var referers = []string{
	"https://www.adidas.com/",
	"https://www.adidas.jp/",
	"https://www.google.com/",
	"https://www.instagram.com/",
}

func SetCommonHeaders(req *http.Request, apiCall bool) {
	rand.Seed(time.Now().UnixNano())
	ref := referers[rand.Intn(len(referers))]
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,ja;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Ch-Ua", "\"Not_A Brand\";v=\"8\", \"Chromium\";v=\"120\", \"Google Chrome\";v=\"120\"")
	req.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	req.Header.Set("Sec-Ch-Ua-Platform", "\"Windows\"")
	if apiCall {
		req.Header.Set("Referer", "https://www.adidas.jp/")
	} else {
		req.Header.Set("Referer", ref)
	}
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

func GetPageHTMLByPassingBot(url string) (string, *html.Node, error) {
	// Create a launcher with leakless disabled
	u := launcher.New().
		Leakless(false). // This disables the use of leakless.exe
		Headless(true). // Run in headless mode (optional)
		MustLaunch()

	// Connect to the launched browser
	browser := rod.New().ControlURL(u).MustConnect()

	// Open the page
	page, err := stealth.Page(browser)
	if err != nil {
		return "", nil, err
	}

	page.Mouse.MustMoveTo(200, 200)
	page.Keyboard.Press(input.Tab)
	time.Sleep(2 * time.Second)

	// Navigate to the provided URL
	err = page.Navigate(url)
	if err != nil {
		return "", nil, err
	}

	// Optional: Wait until the page fully loads
	page.MustWaitLoad()

	//page.MustReload()
	page.MustScreenshotFullPage("")
	// Get HTML
	htmlStr, err := page.HTML()
	if err != nil {
		return "", nil, err
	}

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return "", nil, err
	}

	return htmlStr, doc, nil
}

package common

import (
	"crypto/tls"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

// List of HTTPS proxies (no authentication)
var proxies = []string{
	//"14.241.80.37:443",
	//"43.217.134.23:443",
	//"35.177.23.165:443",
	//"98.130.47.34:443",
	//"3.27.237.252:443",
	//"200.174.198.86:443",
	//"38.147.98.190:443",
	//"8.222.17.214:443",
	//"186.179.169.22:443",
	//"194.170.146.125",
	//"108.136.149.20",

	"85.215.64.49",
}

// NewClientWithProxy creates a new *http.Client using a random proxy
func NewClientWithProxy() *http.Client {
	rand.Seed(time.Now().UnixNano())
	raw := "http://" + proxies[rand.Intn(len(proxies))]

	proxyURL, err := url.Parse(raw)
	if err != nil {
		// Fallback: return default client if parsing fails
		return &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // <--- DISABLE CERT VALIDATION (not secure!)
		},
	}

	return &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}
}

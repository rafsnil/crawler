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
	"38.154.227.167:5868",
	"198.23.239.134:6540",
	"207.244.217.165:6712",
	"107.172.163.27:6543",
	"216.10.27.159",
	"136.0.207.84",
	"64.64.118.149",
	"142.147.128.93",
	"104.239.105.125",
	"206.41.172.74",
}

// NewClientWithProxy creates a new *http.Client using a random proxy
func NewClientWithProxy() *http.Client {
	rand.Seed(time.Now().UnixNano())
	raw := "http://rzbeinzo:00kyb055tmdd@" + proxies[rand.Intn(len(proxies))]

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

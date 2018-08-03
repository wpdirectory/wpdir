package client

import (
	"net"
	"net/http"
	"runtime"
	"time"
)

var apiClient *http.Client
var zipClient *http.Client

// Setup our HTTP transport and client
func init() {
	var netTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}

	// Used for HTTP API Requests
	apiClient = &http.Client{
		Timeout:   time.Second * time.Duration(30),
		Transport: netTransport,
	}

	// Used for downloading Archive files
	zipClient = &http.Client{
		Timeout:   time.Second * time.Duration(180),
		Transport: netTransport,
	}
}

// GetAPI returns an HTTP client
func GetAPI() *http.Client {
	return apiClient
}

// GetZip returns an HTTP client
func GetZip() *http.Client {
	return zipClient
}

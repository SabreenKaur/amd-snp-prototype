package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/bxcodec/httpcache"
	"github.com/google/go-sev-guest/verify/trust"
)

func GetCachingClient() http.Client {
	client := &http.Client{}
	_, err := httpcache.NewWithInmemoryCache(client, true, time.Second*60)
	if err != nil {
		log.Fatal(err)
	}

	return *client
}

type CacheHTTPGetter struct{}

// Define in the header
// func (n *CacheHTTPGetter) Get(url string)

func (n *CacheHTTPGetter) Get(url string) ([]byte, error) {
	client := GetCachingClient()
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("failed to retrieve %s", url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()
	return body, nil
}

func CachingHTTPSGetter() trust.HTTPSGetter {
	return &trust.RetryHTTPSGetter{
		Timeout:       2 * time.Minute,
		MaxRetryDelay: 30 * time.Second,
		Getter:        &CacheHTTPGetter{},
	}
}

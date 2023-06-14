package main

import (
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-sev-guest/verify/trust"
	"github.com/jellydator/ttlcache/v3"
)

type TTLGetter struct {
	manager *Vcek_TTLCache
}

func (n *TTLGetter) Get(url string) ([]byte, error) {
	// Delete all expired items from cache
	n.manager.cache.DeleteExpired()

	// Query the cache for existence of key
	item := n.manager.cache.Get(url, ttlcache.WithDisableTouchOnHit[string, []byte]())

	// Return item if one is found and if is not expired
	if item != nil && !item.IsExpired() {
		fmt.Println("Cache hit!")
		return item.Value(), nil
	}
	// Otherwise query AMD KDS for the certificate
	resp, err := http.Get(url)
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

	// Add CRL to cache and return CRL
	if strings.Contains(url, "crl") {
		fmt.Println("Cache miss in CRL!")
		err := n.addCRLToCache(url, body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return body, nil
	}

	// Add VCEK certificate to cache and return VCEK
	if !strings.Contains(url, "cert_chain") && !strings.Contains(url, "crl") {
		fmt.Println("Cache miss for VCEK!")
		err := n.addVCEKToCache(url, body)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		return body, nil
	}

	// Otherwise just return any other body returned
	return body, nil

}

func (n *TTLGetter) addCRLToCache(url string, body []byte) error {
	crl, err := x509.ParseRevocationList(body)
	if err != nil {
		fmt.Println("Error parsing revocation list")
		return err
	}
	expiry := crl.NextUpdate
	now := time.Now()
	ttl := expiry.Sub(now)
	n.manager.cache.Set(url, body, ttl)
	return nil

}

func (n *TTLGetter) addVCEKToCache(url string, body []byte) error {
	vcekCert, err := x509.ParseCertificate(body)
	if err != nil {
		fmt.Println("Error parsing vcek cert")
		return err
	}
	expiry := vcekCert.NotAfter
	now := time.Now()
	ttl := expiry.Sub(now)
	n.manager.cache.Set(url, body, ttl)
	return nil
}

func TTLHTTPSGetter(manager *Vcek_TTLCache) trust.HTTPSGetter {
	return &trust.RetryHTTPSGetter{
		Timeout:       2 * time.Minute,
		MaxRetryDelay: 30 * time.Second,
		Getter:        &TTLGetter{manager},
	}
}
